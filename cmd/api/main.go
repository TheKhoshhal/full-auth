package main

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"full-auth/internal/handlers"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/css", "assets/css")
	e.Static("/js", "assets/js")

	// Login route
	e.GET("/login", handlers.LoginPage).Name = "userSignInForm"
	e.POST("/login", handlers.Login)

	// Sign up route
	e.GET("/signup", handlers.SignUpPage).Name = "userSignUpForm"
	e.POST("/signup", handlers.SignUp)

	// Unauthenticated route
	e.GET("/", handlers.Accessible)

	// Restricted group
	r := e.Group("/restricted")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecretKey := os.Getenv("JWT_SECRET")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.JwtCustomClaims)
		},
		SigningKey:   []byte(jwtSecretKey),
		TokenLookup:  "cookie:token-access",
		ErrorHandler: handlers.JWTErrorChecker,
	}

	r.Use(echojwt.WithConfig(config))
	r.GET("", handlers.Restricted)

	e.Logger.Fatal(e.Start(":1323"))
}
