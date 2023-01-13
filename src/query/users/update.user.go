package query

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hutanku-service/config"
	"hutanku-service/src/helpers"
	models2 "hutanku-service/src/models"
	"log"
	"os"
	"time"
)

func UpdateUsers(c echo.Context) (models2.Response, error) {
	var res models2.Response
	id := c.QueryParam("id")
	var reqBody models2.Users
	if err := c.Bind(&reqBody); err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	db, err := config.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cancel()

	secret := os.Getenv("RAHASIA_NEGARA")

	NikEncrypted := helper.Encrypt([]byte(reqBody.NIK), secret)
	KkEncrypted := helper.Encrypt([]byte(reqBody.KK), secret)
	PhoneNumberEncrypted := helper.Encrypt([]byte(reqBody.PhoneNumber), secret)
	AlamatEncrypted := helper.Encrypt([]byte(reqBody.Alamat), secret)

	objId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objId}
	update := bson.M{
		"$set": bson.M{
			"fullName":     reqBody.FullName,
			"nomorAnggota": reqBody.NomorAnggota,
			"phoneNumber":  string(PhoneNumberEncrypted),
			"dusun":        reqBody.Dusun,
			"desa":         reqBody.Desa,
			"rt":           reqBody.RT,
			"rw":           reqBody.RW,
			"kecamatan":    reqBody.Kecamatan,
			"kabupaten":    reqBody.Kabupaten,
			"kelurahan":    reqBody.Kelurahan,
			"kota":         reqBody.Kota,
			"alamat":       string(AlamatEncrypted),
			"nik":          string(NikEncrypted),
			"kk":           string(KkEncrypted),
			"province":     reqBody.Province,
			"pokja":        reqBody.Pokja,
		},
	}

	result, err := db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err.Error())
	}

	res.Message = "Updated data success"
	res.Data = result.ModifiedCount

	return res, nil
}
