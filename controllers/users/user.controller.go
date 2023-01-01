package controllers

import (
	"github.com/labstack/echo/v4"
	query "hutanku-service/query/users"
	"net/http"
)

func CreateUsers(c echo.Context) error {
	result, err := query.CreateUsers(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
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

	result, err := query.GetUsers(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func UpdateUsers(c echo.Context) error {
	result, err := query.UpdateUsers(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

func ForgotPasswordUsers(c echo.Context) error {
	result, err := query.ForgotPassword(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func ResetPasswordUser(c echo.Context) error {
	result, err := query.ResetPassword(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
