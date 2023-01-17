package controllers

import (
	"github.com/labstack/echo/v4"
	query "hutanku-service/src/query/lahan"

	"net/http"
)

func GetLahan(c echo.Context) error {

	result, err := query.GetLahanByUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
