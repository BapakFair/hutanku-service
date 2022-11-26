package controllers

import (
	"github.com/labstack/echo/v4"
	"hutanku-service/models"
	query "hutanku-service/query/users"
	"log"
	"net/http"
)

func CreateUsers(c echo.Context) error {
	var reqBody models.Users
	if err := c.Bind(&reqBody); err != nil {
		log.Fatal(err.Error())
	}
	result, err := query.CreateUsers(reqBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

func GetUsers(c echo.Context) error {
	id := c.QueryParam("id")
	if id != "" {
		result, err := query.GetUsersById(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, result)
	}
	result, err := query.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
