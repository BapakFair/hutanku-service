package main

import (
	"github.com/labstack/echo/v4"
	"hutanku-service/config"
	routes "hutanku-service/routes"
	"log"
	"net/http"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
func main() {
	_, err := config.Connect()
	if err != nil {
		log.Fatal(err)
	}

	e := routes.Init()

	e.Logger.Fatal(e.Start(":80"))
}
