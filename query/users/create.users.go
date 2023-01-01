package query

import (
	"context"
	"github.com/labstack/echo/v4"
	"hutanku-service/config"
	helper "hutanku-service/helpers"
	"hutanku-service/models"
	"os"
	"time"
)

func CreateUsers(c echo.Context) (models.Response, error) {
	var res models.Response
	var reqBody models.Users
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
	passwordHashed, _ := helper.HashPassword(reqBody.Password)
	NikEncrypted := helper.Encrypt([]byte(reqBody.NIK), secret)
	KkEncrypted := helper.Encrypt([]byte(reqBody.KK), secret)
	PhoneNumberEncrypted := helper.Encrypt([]byte(reqBody.PhoneNumber), secret)
	AlamatEncrypted := helper.Encrypt([]byte(reqBody.Alamat), secret)

	data, err := db.Collection("users").InsertOne(ctx, models.Users{
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
