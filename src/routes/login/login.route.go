package routes

import (
	"github.com/labstack/echo/v4"
	login "hutanku-service/src/controllers/login"
)

func LoginRoute(e *echo.Echo) {
	e.POST("/api/login", login.Login)
}
