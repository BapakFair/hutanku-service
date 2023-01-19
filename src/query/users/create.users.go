package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hutanku-service/config"
	"hutanku-service/src/helpers"
	"hutanku-service/src/models"
	"os"
	"time"
)

func CreateUsers(c echo.Context) (models.Response, error) {
	var res models.Response
	var reqBodyUser models.Users
	if err := c.Bind(&reqBodyUser); err != nil {
		return res, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	db, err := config.Connect()
	if err != nil {
		return res, err
	}
	defer cancel()

	secret := os.Getenv("RAHASIA_NEGARA")
	passwordHashed, _ := helper.HashPassword(reqBodyUser.Password)
	NikEncrypted := helper.Encrypt([]byte(reqBodyUser.NIK), secret)
	KkEncrypted := helper.Encrypt([]byte(reqBodyUser.KK), secret)
	PhoneNumberEncrypted := helper.Encrypt([]byte(reqBodyUser.PhoneNumber), secret)
	AlamatEncrypted := helper.Encrypt([]byte(reqBodyUser.Alamat), secret)

	// check email exist =============================
	checkEmail, err := db.Collection("users").Find(ctx, bson.M{
		"email": reqBodyUser.Email,
	})
	if err != nil {
		return res, err
	}

	var data []bson.M
	if err := checkEmail.All(ctx, &data); err != nil {
		return res, err
	}
	fmt.Println(data)
	if len(data) != 0 {
		err := errors.New("silahkan gunakan email yang lain, karena email ini sudah ada pemiliknya")
		return res, err
	}
	// ================================================

	dataUser, err := db.Collection("users").InsertOne(ctx, models.Users{
		FullName:     reqBodyUser.FullName,
		NomorAnggota: reqBodyUser.NomorAnggota,
		Email:        reqBodyUser.Email,
		Password:     passwordHashed,
		PhoneNumber:  string(PhoneNumberEncrypted),
		Dusun:        reqBodyUser.Dusun,
		DesaKel:      reqBodyUser.DesaKel,
		RT:           reqBodyUser.RT,
		RW:           reqBodyUser.RW,
		Kecamatan:    reqBodyUser.Kecamatan,
		KotaKab:      reqBodyUser.KotaKab,
		Alamat:       string(AlamatEncrypted),
		NIK:          string(NikEncrypted),
		KK:           string(KkEncrypted),
		Province:     reqBodyUser.Province,
		Pokja:        reqBodyUser.Pokja,
		Role:         reqBodyUser.Role,
	})
	if err != nil {
		return res, err
	}

	if len(reqBodyUser.DataPetak) != 0 {
		for _, dataPetak := range reqBodyUser.DataPetak {
			checkPetak, err := db.Collection("petak").Find(ctx, bson.M{
				"petak": dataPetak["petak"],
				"andil": dataPetak["andil"],
			})
			if err != nil {
				return res, err
			}

			var data []bson.M
			if err := checkPetak.All(ctx, &data); err != nil {
				return res, err
			}
			fmt.Println(data)
			if len(data) != 0 {
				err := errors.New("silahkan pilih petak yang lain, karena petak ini sudah ada pemiliknya")
				return res, err
			}

			_, err = db.Collection("petak").InsertOne(ctx, models.Petak{
				UserId:    dataUser.InsertedID.(primitive.ObjectID),
				Petak:     dataPetak["petak"].(string),
				Andil:     dataPetak["andil"].(string),
				Pokja:     dataPetak["pokja"].(string),
				LuasLahan: dataPetak["luasLahan"].(float64),
			})
			if err != nil {
				return res, err
			}
		}
	}

	res.Message = "Insert data success"
	res.Data = dataUser

	return res, nil
}
