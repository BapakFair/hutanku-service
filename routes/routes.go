package routes

import (
	"github.com/labstack/echo/v4"
	login "hutanku-service/controllers/login"
	user "hutanku-service/controllers/users"
	"hutanku-service/middleware"
	"net/http"
)

func Init() *echo.Echo {

	e := echo.New()

	e.Use(middleware.Logger)
	e.Use(middleware.Recover)
	e.Use(middleware.MidCors)
	
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hallo this is hutanku")
	})
	e.GET("/api/users", user.GetUsers, middleware.IsAuthenticated)
	e.POST("/api/user", user.CreateUsers, middleware.IsAuthenticated)
	e.PUT("/api/user/update", user.UpdateUsers)
	e.POST("/api/user/forgot", user.ForgotPasswordUsers)
	e.POST("/api/user/reset", user.ResetPasswordUser)

	e.POST("/api/login", login.Login)

	return e
}
