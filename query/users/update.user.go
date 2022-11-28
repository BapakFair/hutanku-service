package query

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hutanku-service/config"
	helper "hutanku-service/helpers"
	"hutanku-service/models"
	"log"
	"net/http"
	"os"
	"time"
)

func UpdateUsers(id string, reqBody models.Users) (models.Response, error) {
	var res models.Response

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
			"FullName":     reqBody.FullName,
			"NomorAnggota": reqBody.NomorAnggota,
			"PhoneNumber":  string(PhoneNumberEncrypted),
			"Dusun":        reqBody.Dusun,
			"Desa":         reqBody.Desa,
			"RT":           reqBody.RT,
			"RW":           reqBody.RW,
			"Kecamatan":    reqBody.Kecamatan,
			"Kabupaten":    reqBody.Kabupaten,
			"Kelurahan":    reqBody.Kelurahan,
			"Kota":         reqBody.Kota,
			"Alamat":       string(AlamatEncrypted),
			"NIK":          string(NikEncrypted),
			"KK":           string(KkEncrypted),
			"Province":     reqBody.Province,
			"Pokja":        reqBody.Pokja,
		},
	}

	result, err := db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err.Error())
	}

	res.Status = http.StatusCreated
	res.Message = "Insert data success"
	res.Data = result

	return res, nil
}
