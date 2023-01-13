package middleware

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

var envErr = godotenv.Load()
var kunci string = os.Getenv("KUNCI_MASUK")

var IsAuthenticated = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(kunci),
})

//var MidCors = middleware.CORSWithConfig(middleware.CORSConfig{
//	AllowOrigins: []string{"*"},
//	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
//})

var MidCors = middleware.CORS()

var Logger = middleware.Logger()

var Recover = middleware.Recover()
