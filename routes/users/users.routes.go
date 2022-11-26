package routes

import (
	"github.com/labstack/echo/v4"
	controllers "hutanku-service/controllers/users"
	"net/http"
)

func Init() *echo.Echo {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hallo this is hutanku")
	})
	e.GET("/api/users", controllers.GetUsers)
	e.POST("/api/user", controllers.CreateUsers)

	return e
}
