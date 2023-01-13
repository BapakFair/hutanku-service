package routes

import (
	"github.com/labstack/echo/v4"
	dashboard "hutanku-service/src/controllers/dashboard"
	"hutanku-service/src/middleware"
)

func DashboardAdminRoutes(e *echo.Echo) {
	e.GET("/api/dashboard", dashboard.GetDashboardData, middleware.IsAuthenticated)
}
