package routes

import (
	"github.com/labstack/echo/v4"
	"hutanku-service/src/controllers/lahan"
	"hutanku-service/src/middleware"
)

func LahanRoutes(e *echo.Echo) {
	e.GET("/api/lahan", controllers.GetLahan, middleware.IsAuthenticated)
}
