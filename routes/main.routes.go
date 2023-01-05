package routes

import (
	"github.com/labstack/echo/v4"
	"hutanku-service/middleware"
	routesAdm "hutanku-service/routes/adm"
	routesDashboard "hutanku-service/routes/dashboard"
	routesLogin "hutanku-service/routes/login"
	routesUsers "hutanku-service/routes/users"

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
	routesUsers.UserRoutes(e)
	routesLogin.LoginRoute(e)
	routesAdm.AdmRoutes(e)
	routesDashboard.DashboardAdminRoutes(e)

	return e
}
