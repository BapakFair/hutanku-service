package query

import (
	"context"
	"github.com/labstack/echo/v4"
	"hutanku-service/config"
	helper2 "hutanku-service/src/helpers"
	models2 "hutanku-service/src/models"
	"os"
	"time"
)

func CreateUsers(c echo.Context) (models2.Response, error) {
	var res models2.Response
	var reqBody models2.Users
	if err := c.Bind(&reqBody); err != nil {
		return res, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	db, err := config.Connect()
	if err != nil {
		return res, err
	}
	defer cancel()

	secret := os.Getenv("RAHASIA_NEGARA")
	passwordHashed, _ := helper2.HashPassword(reqBody.Password)
	NikEncrypted := helper2.Encrypt([]byte(reqBody.NIK), secret)
	KkEncrypted := helper2.Encrypt([]byte(reqBody.KK), secret)
	PhoneNumberEncrypted := helper2.Encrypt([]byte(reqBody.PhoneNumber), secret)
	AlamatEncrypted := helper2.Encrypt([]byte(reqBody.Alamat), secret)

	data, err := db.Collection("users").InsertOne(ctx, models2.Users{
		FullName:     reqBody.FullName,
		NomorAnggota: reqBody.NomorAnggota,
		Email:        reqBody.Email,
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
		return res, err
	}

	res.Message = "Insert data success"
	res.Data = data

	return res, nil
}
