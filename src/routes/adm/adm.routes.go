package routes

import (
	"github.com/labstack/echo/v4"
	adm "hutanku-service/src/controllers/adm"
	"hutanku-service/src/middleware"
)

func AdmRoutes(e *echo.Echo) {
	e.GET("/api/adm", adm.GetAdm, middleware.IsAuthenticated)
}
