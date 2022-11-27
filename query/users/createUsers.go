package query

import (
	"context"
	"hutanku-service/config"
	helper "hutanku-service/helpers"
	"hutanku-service/models"
	"log"
	"net/http"
	"os"
	"time"
)

func CreateUsers(reqBody models.Users) (models.Response, error) {
	var res models.Response
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	db, err := config.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cancel()

	secret := os.Getenv("RAHASIA_NEGARA")
	passwordHashed, _ := helper.HashPassword(reqBody.Password)
	NikEncrypted := helper.Encrypt([]byte(reqBody.NIK), secret)
	KkEncrypted := helper.Encrypt([]byte(reqBody.KK), secret)
	PhoneNumberEncrypted := helper.Encrypt([]byte(reqBody.PhoneNumber), secret)
	AlamatEncrypted := helper.Encrypt([]byte(reqBody.Alamat), secret)
	EmailEncrypted := helper.Encrypt([]byte(reqBody.Email), secret)

	data, err := db.Collection("users").InsertOne(ctx, models.Users{
		FullName:     reqBody.FullName,
		NomorAnggota: reqBody.NomorAnggota,
		Email:        string(EmailEncrypted),
		Password:     passwordHashed,
		PhoneNumber:  string(PhoneNumberEncrypted),
		Dusun:        reqBody.Dusun,
		Desa:         reqBody.Desa,
		RT:           reqBody.RT,
		RW:           reqBody.RW,
		Kecamatan:    reqBody.Kecamatan,
		Kabupaten:    reqBody.Kabupaten,
		Kelurahan:    reqBody.Kelurahan,
		Kota:         reqBody.Kota,
		Alamat:       string(AlamatEncrypted),
		NIK:          string(NikEncrypted),
		KK:           string(KkEncrypted),
		Province:     reqBody.Province,
		Pokja:        reqBody.Pokja,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	res.Status = http.StatusCreated
	res.Message = "Insert data success"
	res.Data = data

	return res, nil
}
