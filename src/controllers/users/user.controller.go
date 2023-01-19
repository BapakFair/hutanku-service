package controllers

import (
	"github.com/labstack/echo/v4"
	query2 "hutanku-service/src/query/users"
	"net/http"
)

func CreateUsers(c echo.Context) error {
	result, err := query2.CreateUsers(c)
	if err.Data != nil {
		return c.JSON(err.Status, map[string]string{"message": err.Data.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

func GetUsers(c echo.Context) error {
	// get data from token
	//userData := c.Get("user").(*jwt.Token)
	//claims := userData.Claims.(jwt.MapClaims)
	//role := claims["role"].(float64)
	//noAnggota := claims["nomorAnggota"].(float64)
	// =========================================

	result, err := query2.GetUsers(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func UpdateUsers(c echo.Context) error {
	result, err := query2.UpdateUsers(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

func ForgotPasswordUsers(c echo.Context) error {
	result, err := query2.ForgotPassword(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func ResetPasswordUser(c echo.Context) error {
	result, err := query2.ResetPassword(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
