package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// @todo #1:15m Write some documentation about JWT middleware usage
func GetJWTMiddleware(signingKey string) echo.MiddlewareFunc {
	return middleware.JWT([]byte(signingKey))
}
