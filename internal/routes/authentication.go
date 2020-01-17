package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

func begin(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func verify(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func refresh(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
