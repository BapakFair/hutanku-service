package query

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hutanku-service/config"
	"hutanku-service/src/models"
	"time"
)

func GetLahanByUser(c echo.Context) (models.Response, error) {
	var res models.Response
	id, _ := primitive.ObjectIDFromHex(c.QueryParam("id"))
	var data []bson.M

	ctx, cancel := context.WithTimeout(context.Background(), 40000*time.Second)
	db, err := config.Connect()
	if err != nil {
		return res, err
	}

	defer cancel()

	dataLahan, err := db.Collection("petak").Find(ctx, bson.M{"userId": id})

	if err != nil {
		return res, err
	}
	if err := dataLahan.All(ctx, &data); err != nil {
		return res, err
	}

	if len(data) == 0 {
		err := errors.New("no documents in result")
		return res, err
	}

	res.Message = "Get data success"
	res.Data = data

	return res, nil
}
