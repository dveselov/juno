package routes

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/dveselov/juno/internal/authentication/providers"
	"github.com/dveselov/juno/internal/config"
	"github.com/dveselov/juno/internal/utils"
)

type AuthenticationBeginRequest struct {
	PhoneNumber string `json:"phoneNumber"`
}

func (r *AuthenticationBeginRequest) Validate() error {
	return utils.IsInvalidPhoneNumber(r.PhoneNumber)
}

type AuthenticationBeginResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func Begin(c echo.Context) error {
	request := new(AuthenticationBeginRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, AuthenticationBeginResponse{
			Status:  false,
			Message: err.Error(),
		})
	}
	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, AuthenticationBeginResponse{
			Status:  false,
			Message: err.Error(),
		})
	}
	code := utils.GetRandDigitsString(4)
	provider := c.(*config.AppContext).GetAuthenticationProvider()
	err := providers.Begin(provider, request.PhoneNumber, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthenticationBeginResponse{
			Status:  false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, AuthenticationBeginResponse{
		Status:  true,
		Message: "OK",
	})
}

type AuthenticationVerifyRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Code        string `json:"code"`
}

func (r *AuthenticationVerifyRequest) Validate() error {
	if err := utils.IsInvalidPhoneNumber(r.PhoneNumber); err != nil {
		return err
	}
	if err := utils.IsInvalidVerificationCode(r.Code); err != nil {
		return err
	}
	return nil
}

type AuthenticationVerifyResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
}

func Verify(c echo.Context) error {
	request := new(AuthenticationVerifyRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, AuthenticationVerifyResponse{
			Status:  false,
			Message: err.Error(),
		})
	}
	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, AuthenticationVerifyResponse{
			Status:  false,
			Message: err.Error(),
		})
	}
	context := c.(*config.AppContext)

	provider := context.GetAuthenticationProvider()
	err := providers.Verify(provider, request.PhoneNumber, request.Code)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, AuthenticationVerifyResponse{
			Status:  false,
			Message: "Invalid phone number or verification code",
		})
	}

	db := context.GetDB()

	stmt, err := db.Prepare(`
		INSERT INTO authentication_user(phone_number) VALUES ($1)
		ON CONFLICT(phone_number)
		DO UPDATE
		SET phone_number = $1
		RETURNING id;
	`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthenticationBeginResponse{
			Status:  false,
			Message: err.Error(),
		})
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthenticationBeginResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	var userID int

	err = stmt.QueryRow(request.PhoneNumber).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, AuthenticationBeginResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthenticationBeginResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = string(userID)
	claims["phone_number"] = request.PhoneNumber
	// @todo #1:30min Use token expiration time from env variables
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(context.Config.SigningKey))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthenticationBeginResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, AuthenticationVerifyResponse{
		Status:  true,
		Message: "OK",
		Token:   t,
	})
}
