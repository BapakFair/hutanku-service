package query

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"hutanku-service/config"
	helper "hutanku-service/helpers"
	"hutanku-service/models"
	"log"
	"net/http"
	"time"
)

func ForgotPassword(reqBody models.Users) (models.Response, error) {
	var res models.Response
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	db, err := config.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cancel()

	passwordHashed, _ := helper.HashPassword(reqBody.Email + "secretTokenForResetPassword")

	filter := bson.M{"email": reqBody.Email}
	update := bson.M{
		"$set": bson.M{
			"resetToken": passwordHashed,
		},
	}

	result, err := db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err.Error())
	}

	res.Status = http.StatusOK
	res.Message = "Please Check Your Email"
	res.Data = result.UpsertedID

	return res, nil
}
