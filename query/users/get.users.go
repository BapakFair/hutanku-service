package query

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hutanku-service/config"
	helper "hutanku-service/helpers"
	"hutanku-service/models"
	"net/http"
	"os"
	"time"
)

func GetUsersById(id string) (models.Response, error) {
	var res models.Response
	var users models.Users
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := config.Connect()
	if err != nil {
		return res, err
	}

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)

	if err := db.Collection("users").FindOne(ctx, bson.M{"_id": objId}).Decode(&users); err != nil {
		return res, err
	}
	secret := os.Getenv("RAHASIA_NEGARA")

	users.NIK = string(helper.Decrypt([]byte(users.NIK), secret))
	users.KK = string(helper.Decrypt([]byte(users.KK), secret))
	users.PhoneNumber = string(helper.Decrypt([]byte(users.PhoneNumber), secret))
	users.Alamat = string(helper.Decrypt([]byte(users.Alamat), secret))

	res.Status = http.StatusOK
	res.Message = "Get data success"
	res.Data = users

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
	if err := data.All(ctx, &dataFinal); err != nil {
		return res, err
	}

	if len(dataFinal) == 0 {
		err = errors.New("no documents in result")
		return res, err
	}

	for i := 0; i < len(dataFinal); i++ {
		secret := os.Getenv("RAHASIA_NEGARA")
		NIK := fmt.Sprintf("%v", dataFinal[i]["nik"])
		KK := fmt.Sprintf("%v", dataFinal[i]["kk"])
		Phone := fmt.Sprintf("%v", dataFinal[i]["phoneNumber"])
		Alamat := fmt.Sprintf("%v", dataFinal[i]["alamat"])

		dataFinal[i]["nik"] = string(helper.Decrypt([]byte(NIK), secret))
		dataFinal[i]["kk"] = string(helper.Decrypt([]byte(KK), secret))
		dataFinal[i]["phoneNumber"] = string(helper.Decrypt([]byte(Phone), secret))
		dataFinal[i]["alamat"] = string(helper.Decrypt([]byte(Alamat), secret))
	}

	res.Status = http.StatusOK
	res.Message = "Get data success"
	res.Data = dataFinal

	return res, nil
}
