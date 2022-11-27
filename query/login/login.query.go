package query

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"hutanku-service/config"
	helper "hutanku-service/helpers"
	"hutanku-service/models"
	"log"
	"time"
)

func CheckLogin(email, password string) (bool, error) {

	var user models.Users

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	db, err := config.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cancel()

	if err := db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return false, err
	}

	match, err := helper.CheckPasswordHash(password, user.Password)
	if !match {
		errors.New("email or password keliru bosku ...")
		return false, err
	}

	return true, nil
}
