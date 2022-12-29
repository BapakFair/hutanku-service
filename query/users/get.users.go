package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hutanku-service/config"
	helper "hutanku-service/helpers"
	"hutanku-service/models"
	"os"
	"strconv"
	"time"
)

func GetUsersById(id string) (models.Response, error) {
	var res models.Response
	var users models.Users
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := config.Connect()
	if err != nil {
		return res, err
	}

	objId, _ := primitive.ObjectIDFromHex(id)

	if err := db.Collection("users").FindOne(ctx, bson.M{"_id": objId}).Decode(&users); err != nil {
		return res, err
	}
	secret := os.Getenv("RAHASIA_NEGARA")

	users.NIK = string(helper.Decrypt([]byte(users.NIK), secret))
	users.KK = string(helper.Decrypt([]byte(users.KK), secret))
	users.PhoneNumber = string(helper.Decrypt([]byte(users.PhoneNumber), secret))
	users.Alamat = string(helper.Decrypt([]byte(users.Alamat), secret))

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

	res.Message = "Get data success"
	res.Data = dataFinal

	return res, nil
}

func GetUserByNoAnggota(nomorAnggota string) (models.Response, error) {
	var res models.Response
	var users models.Users
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := config.Connect()
	if err != nil {
		return res, err
	}

	defer cancel()
	nomorAnggotaInt, _ := strconv.Atoi(nomorAnggota)
	if err := db.Collection("users").FindOne(ctx, bson.M{"nomorAnggota": nomorAnggotaInt}).Decode(&users); err != nil {
		return res, err
	}
	secret := os.Getenv("RAHASIA_NEGARA")

	users.NIK = string(helper.Decrypt([]byte(users.NIK), secret))
	users.KK = string(helper.Decrypt([]byte(users.KK), secret))
	users.PhoneNumber = string(helper.Decrypt([]byte(users.PhoneNumber), secret))
	users.Alamat = string(helper.Decrypt([]byte(users.Alamat), secret))

	res.Message = "Get data success"
	res.Data = users

	return res, nil
}

func GetUsers(c echo.Context) (models.Response, error) {
	var res models.Response
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	db, err := config.Connect()
	if err != nil {
		return res, err
	}

	defer cancel()

	filter := bson.M{}
	findOptions := options.Find()
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("perPage"))
	findOptions.SetSkip((int64(page) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))
	//findOptions.SetProjection(bson.M{"_id": 1})

	if q := c.QueryParam("q"); q != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					//	"nik": bson.M{
					//		"$regex": primitive.Regex{
					//			Pattern: q,
					//			Options: "i",
					//		},
					//	},
					//}, {
					//	"kk": bson.M{
					//		"$regex": primitive.Regex{
					//			Pattern: q,
					//			Options: "i",
					//		},
					//	},
					//}, {
					//	"nomorAnggota": bson.M{
					//		"$regex": primitive.Regex{
					//			Pattern: q,
					//			Options: "i",
					//		},
					//	},
					//}, {
					"fullName": bson.M{
						"$regex": primitive.Regex{
							Pattern: q,
							Options: "i",
						},
					},
				},
			},
		}
	}
	data, err := db.Collection("users").Find(ctx, filter, findOptions)
	if err != nil {
		return res, err
	}
	var dataFinal []bson.M
	if err := data.All(ctx, &dataFinal); err != nil {
		return res, err
	}
	fmt.Println("test data : ", len(dataFinal))
	if len(dataFinal) == 0 {
		err = errors.New("no documents in result")
		return res, err
	}

	// this line of code below to manual hash nik & kk data from string to hashed =======================
	// don't forget to change context timeout 120 second per 1000 data
	//secret := os.Getenv("RAHASIA_NEGARA")
	//for i := 0; i < len(dataFinal); i++ {
	//	objId, _ := dataFinal[i]["_id"].(primitive.ObjectID)
	//	findOption := options.Find()
	//	findOption.SetProjection(bson.M{
	//		"nik": 1,
	//		"kk":  1,
	//	})
	//	var dataReadyUpdate []bson.M
	//	dataPerUser, _ := db.Collection("users").Find(ctx, bson.M{"_id": objId}, findOption)
	//	if err := dataPerUser.All(ctx, &dataReadyUpdate); err != nil {
	//		return res, err
	//	}
	//	dataNik := dataReadyUpdate[0]["nik"].(string)
	//	dataKk := dataReadyUpdate[0]["kk"].(string)
	//	if len(dataNik) == 16 && len(dataKk) == 16 {
	//		NIKhashed := helper.Encrypt([]byte(dataNik), secret)
	//		KKhashed := helper.Encrypt([]byte(dataKk), secret)
	//		filter := bson.M{"_id": objId}
	//		update := bson.M{
	//			"$set": bson.M{
	//				"nik": string(NIKhashed),
	//				"kk":  string(KKhashed),
	//			},
	//		}
	//		db.Collection("users").UpdateOne(ctx, filter, update)
	//	}
	// =============================================================================================
	//fmt.Println("ini data test :", "data ke :", i)
	for i := 0; i < len(dataFinal); i++ {
		secret := os.Getenv("RAHASIA_NEGARA")
		NIK := fmt.Sprintf("%v", dataFinal[i]["nik"])
		KK := fmt.Sprintf("%v", dataFinal[i]["kk"])
		//Phone := fmt.Sprintf("%v", dataFinal[i]["phoneNumber"])
		//Alamat := fmt.Sprintf("%v", dataFinal[i]["alamat"])

		dataFinal[i]["nik"] = string(helper.Decrypt([]byte(NIK), secret))
		dataFinal[i]["kk"] = string(helper.Decrypt([]byte(KK), secret))
		//dataFinal[i]["phoneNumber"] = string(helper.Decrypt([]byte(Phone), secret))
		//dataFinal[i]["alamat"] = string(helper.Decrypt([]byte(Alamat), secret))
	}
	res.Message = "Get data success"
	res.Data = dataFinal

	return res, nil
}
