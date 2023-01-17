package query

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hutanku-service/config"
	"hutanku-service/src/helpers"
	"hutanku-service/src/models"
	"log"
	"os"
	"time"
)

func UpdateUsers(c echo.Context) (models.Response, error) {
	var res models.Response
	id := c.QueryParam("id")
	var reqBody models.Users
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
			"desakel":      reqBody.DesaKel,
			"rt":           reqBody.RT,
			"rw":           reqBody.RW,
			"kecamatan":    reqBody.Kecamatan,
			"kotakab":      reqBody.KotaKab,
			"alamat":       string(AlamatEncrypted),
			"nik":          string(NikEncrypted),
			"kk":           string(KkEncrypted),
			"province":     reqBody.Province,
			"pokja":        reqBody.Pokja,
			"role":         reqBody.Role,
		},
	}

	result, err := db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(reqBody.DataPetak) != 0 {
		for _, dataPetak := range reqBody.DataPetak {
			checkPetak, err := db.Collection("petak").Find(ctx, bson.M{
				"petak": dataPetak["petak"],
				"andil": dataPetak["andil"],
			})
			if err != nil {
				return res, err
			}
			if checkPetak != nil {
				err := errors.New("silahkan pilih petak dan andil yang lain, karena petak ini sudah ada pemiliknya")
				return res, err
			}

			filterLahan := bson.M{"userId": objId}
			updateLahan := bson.M{
				"$set": bson.M{
					"petak":     dataPetak["petak"].(string),
					"andil":     dataPetak["andil"].(string),
					"pokja":     dataPetak["pokja"].(string),
					"luasLahan": dataPetak["luasLahan"].(float64),
				},
			}
			_, err = db.Collection("petak").UpdateOne(ctx, filterLahan, updateLahan)
			if err != nil {
				return res, err
			}
		}
	}

	res.Message = "Updated data success"
	res.Data = result.ModifiedCount

	return res, nil
}
