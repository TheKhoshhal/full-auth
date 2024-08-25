package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"full-auth/cmd/web"
	"full-auth/internal/db_setup"
	"full-auth/internal/renderer"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	usernameDb, passwordDb := db_setup.GetUser(username)

	// Throws unauthorized error
	if username != usernameDb {
		return renderer.Render(c, http.StatusOK, web.ShowSignIn(true))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordDb), []byte(password)); err != nil {
		return renderer.Render(c, http.StatusOK, web.ShowSignIn(true))
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		username,
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecretKey := os.Getenv("JWT_SECRET")

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return err
	}

	setTokenCookie(t, c)

	return c.Redirect(http.StatusMovedPermanently, "/restricted")
}

func SignUp(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	usernameDb, _ := db_setup.GetUser(username)
	if usernameDb != "" {
		return renderer.Render(c, http.StatusOK, web.ShowSignUp(true))
	}

	accountName := db_setup.CreateAccount(username, string(hashedPassword), email)

	// Set custom claims
	claims := &JwtCustomClaims{
		accountName,
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecretKey := os.Getenv("JWT_SECRET")

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return err
	}

	setTokenCookie(t, c)

	return c.Redirect(http.StatusMovedPermanently, "/restricted")
}

// func DeleteAccount(c echo.Context) error {
// 	claims := c.Get("user").(*jwt.Token)
// 	name := claims.Claims.(*JwtCustomClaims).Name
//
// 	db_setup.DeleteAccount(name)
//
// 	return c.Redirect(http.StatusMovedPermanently, "/")
// }

func setTokenCookie(token string, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "token-access"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * 72)
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func Accessible(c echo.Context) error {
	return renderer.Render(c, http.StatusOK, web.ShowAccessible())
}

func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	name := claims.Name
	return renderer.Render(c, http.StatusOK, web.ShowRestricted(name))
}

func LoginPage(c echo.Context) error {
	return renderer.Render(c, http.StatusOK, web.ShowSignIn(false))
}

func SignUpPage(c echo.Context) error {
	return renderer.Render(c, http.StatusOK, web.ShowSignUp(false))
}
