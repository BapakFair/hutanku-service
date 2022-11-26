package query

import (
	"context"
	"hutanku-service/config"
	"hutanku-service/models"
	"log"
	"net/http"
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

	data, err := db.Collection("users").InsertOne(ctx, models.Users{
		FullName:     reqBody.FullName,
		NomorAnggota: reqBody.NomorAnggota,
		Email:        reqBody.Email,
		Password:     reqBody.Password,
		PhoneNumber:  reqBody.PhoneNumber,
		Dusun:        reqBody.Dusun,
		Desa:         reqBody.Desa,
		RT:           reqBody.RT,
		RW:           reqBody.RW,
		Kecamatan:    reqBody.Kecamatan,
		Kabupaten:    reqBody.Kabupaten,
		Kelurahan:    reqBody.Kelurahan,
		Kota:         reqBody.Kota,
		Alamat:       reqBody.Alamat,
		NIK:          reqBody.NIK,
		KK:           reqBody.KK,
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
