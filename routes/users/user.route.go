package routes

import (
	"github.com/labstack/echo/v4"
	user "hutanku-service/controllers/users"
	"hutanku-service/middleware"
)

func UserRoutes(e *echo.Echo) {
	e.GET("/api/users", user.GetUsers, middleware.IsAuthenticated)
	e.POST("/api/user", user.CreateUsers, middleware.IsAuthenticated)
	e.PUT("/api/user/update", user.UpdateUsers)
	e.POST("/api/user/forgot", user.ForgotPasswordUsers)
	e.POST("/api/user/reset", user.ResetPasswordUser)
}
