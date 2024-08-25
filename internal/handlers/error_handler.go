package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func JWTErrorChecker(c echo.Context, err error) error {
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
}
