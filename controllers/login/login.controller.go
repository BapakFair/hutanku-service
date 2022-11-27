package controllers

import (
	"github.com/labstack/echo/v4"
	"hutanku-service/models"
	query "hutanku-service/query/login"
	"log"
	"net/http"
)

func Login(c echo.Context) error {
	var reqBody models.Login
	if err := c.Bind(&reqBody); err != nil {
		log.Fatal(err.Error())
	}

	result, err := query.CheckLogin(reqBody.Email, reqBody.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
