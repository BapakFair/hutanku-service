package query

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hutanku-service/config"
	"hutanku-service/models"
	"net/http"
	"time"
)

func GetUsersById(id string) (models.Response, error) {
	var res models.Response
	var product models.Users
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := config.Connect()
	if err != nil {
		return res, err
	}

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)

	if err := db.Collection("users").FindOne(ctx, bson.M{"_id": objId}).Decode(&product); err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Get data success"
	res.Data = product

	return res, nil
}

func GetAllUsers() (models.Response, error) {
	var res models.Response
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := config.Connect()
	if err != nil {
		return res, err
	}

	defer cancel()

	data, err := db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return res, err
	}
	var dataFinal []bson.M
	if err = data.All(ctx, &dataFinal); err != nil {
		return res, err
	}
	res.Status = http.StatusOK
	res.Message = "Get data success"
	res.Data = dataFinal

	return res, nil
}
