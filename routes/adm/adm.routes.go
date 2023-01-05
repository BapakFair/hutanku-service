package routes

import (
	"github.com/labstack/echo/v4"
	adm "hutanku-service/controllers/adm"
	"hutanku-service/middleware"
)

func AdmRoutes(e *echo.Echo) {
	e.GET("/api/adm", adm.GetAdm, middleware.IsAuthenticated)
}
