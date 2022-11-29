package controllers

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"hutanku-service/models"
	query "hutanku-service/query/users"
	"log"
	"net/http"
	"strconv"
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
	// get data from token
	userData := c.Get("user").(*jwt.Token)
	claims := userData.Claims.(jwt.MapClaims)
	role := claims["role"].(float64)
	noAnggota := claims["nomorAnggota"].(float64)
	// =========================================
	fmt.Println(role)
	id := c.QueryParam("id")
	nomorAnggota := c.QueryParam("na")

	if id != "" {
		result, err := query.GetUsersById(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, result)
	}
	if role == 0 && nomorAnggota == "" && id == "" {
		result, err := query.GetAllUsers()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusOK, result)
	} else if (nomorAnggota != "" && role == 0) || (nomorAnggota != "" && strconv.Itoa(int(noAnggota)) == nomorAnggota) {
		result, err := query.GetUserByNoAnggota(nomorAnggota)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, result)
	} else {
		err := errors.New("only Admin can see this data")
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
	}
}

func UpdateUsers(c echo.Context) error {
	id := c.QueryParam("id")
	var reqBody models.Users
	if err := c.Bind(&reqBody); err != nil {
		log.Fatal(err.Error())
	}
	result, err := query.UpdateUsers(id, reqBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

func ForgotPasswordUsers(c echo.Context) error {
	var reqBody models.Users
	if err := c.Bind(&reqBody); err != nil {
		log.Fatal(err.Error())
	}
	result, err := query.ForgotPassword(reqBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func ResetPasswordUser(c echo.Context) error {
	var reqBody models.ResetPassword
	if err := c.Bind(&reqBody); err != nil {
		log.Fatal(err.Error())
	}
	result, err := query.ResetPassword(reqBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
