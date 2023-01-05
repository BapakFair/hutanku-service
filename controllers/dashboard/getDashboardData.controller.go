package controllers

import (
	"github.com/labstack/echo/v4"
	query "hutanku-service/query/dashboard"
	"net/http"
)

func GetDashboardData(c echo.Context) error {
	result, err := query.GetHeaderDashboardData(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
