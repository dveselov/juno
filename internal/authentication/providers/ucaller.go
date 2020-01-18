package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type UCallerProvider struct {
	Database   *sqlx.DB
	ServiceID  string
	SecretKey  string
	BaseAPIUrl string
}

func (p UCallerProvider) GetDB() *sqlx.DB {
	return p.Database
}

func (p UCallerProvider) GetDBType() string {
	return "ucaller"
}

func (p UCallerProvider) Send(phoneNumber string, code string) error {
	request, err := http.NewRequest("GET", p.BaseAPIUrl+"/initCall", nil)
	if err != nil {
		return err
	}
	request.Header.Add("Accept", "application/json")

	query := request.URL.Query()
	query.Add("service_id", p.ServiceID)
	query.Add("key", p.SecretKey)
	query.Add("phone", phoneNumber)
	query.Add("code", code)
	request.URL.RawQuery = query.Encode()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		message := fmt.Sprintf("Invalid status code: %d", response.StatusCode)
		return errors.New(message)
	}
	type Response struct {
		Status bool   `json:"status"`
		Error  string `json:"error"`
	}
	body := Response{}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return err
	}
	if !body.Status {
		return errors.New(body.Error)
	}
	return nil
}
