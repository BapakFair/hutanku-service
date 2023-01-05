package routes

import (
	"github.com/labstack/echo/v4"
	dashboard "hutanku-service/controllers/dashboard"
	"hutanku-service/middleware"
)

func DashboardAdminRoutes(e *echo.Echo) {
	e.GET("/api/dashboard", dashboard.GetDashboardData, middleware.IsAuthenticated)
}
