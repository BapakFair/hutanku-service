package controllers

import (
	"github.com/labstack/echo/v4"
	query "hutanku-service/query/adm"
	"net/http"
)

func GetAdm(c echo.Context) error {
	result, err := query.GetAdm(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
