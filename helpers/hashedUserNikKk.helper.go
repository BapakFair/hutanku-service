package helper

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hutanku-service/config"
	"os"
)

func EncryptNikKk(dataFinal []bson.M, ctx context.Context) error {
	db, err := config.Connect()
	if err != nil {
		return err
	}
	secret := os.Getenv("RAHASIA_NEGARA")
	for i := 0; i < len(dataFinal); i++ {
		objId, _ := dataFinal[i]["_id"].(primitive.ObjectID)
		findOption := options.Find()
		findOption.SetProjection(bson.M{
			"nik": 1,
			"kk":  1,
		})
		var dataReadyUpdate []bson.M
		dataPerUser, _ := db.Collection("users").Find(ctx, bson.M{"_id": objId}, findOption)
		if err := dataPerUser.All(ctx, &dataReadyUpdate); err != nil {
			return err
		}
		dataNik := dataReadyUpdate[0]["nik"].(string)
		dataKk := dataReadyUpdate[0]["kk"].(string)
		if len(dataNik) == 16 && len(dataKk) == 16 {
			NIKhashed := Encrypt([]byte(dataNik), secret)
			KKhashed := Encrypt([]byte(dataKk), secret)
			filter := bson.M{"_id": objId}
			update := bson.M{
				"$set": bson.M{
					"nik": string(NIKhashed),
					"kk":  string(KKhashed),
				},
			}
			_, err := db.Collection("users").UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DecryptNikKk(dataFinal []bson.M, ctx context.Context) error {
	db, err := config.Connect()
	if err != nil {
		return err
	}
	secret := os.Getenv("RAHASIA_NEGARA")
	for i := 0; i < len(dataFinal); i++ {
		objId, _ := dataFinal[i]["_id"].(primitive.ObjectID)
		findOption := options.Find()
		findOption.SetProjection(bson.M{
			"nik": 1,
			"kk":  1,
		})
		var dataReadyUpdate []bson.M
		dataPerUser, _ := db.Collection("users").Find(ctx, bson.M{"_id": objId}, findOption)
		if err := dataPerUser.All(ctx, &dataReadyUpdate); err != nil {
			return err
		}
		dataNik := dataReadyUpdate[0]["nik"].(string)
		dataKk := dataReadyUpdate[0]["kk"].(string)

		NIKhashed := Decrypt([]byte(dataNik), secret)
		KKhashed := Decrypt([]byte(dataKk), secret)
		filter := bson.M{"_id": objId}
		update := bson.M{
			"$set": bson.M{
				"nik": string(NIKhashed),
				"kk":  string(KKhashed),
			},
		}
		_, err := db.Collection("users").UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
		fmt.Println("data ke : ", i)

	}
	return nil
}

func UpdateIdUserToPetak(dataFinal []bson.M, ctx context.Context) error {
	db, err := config.Connect()
	if err != nil {
		return err
	}

	for i := 0; i < len(dataFinal); i++ {
		//objId, _ := dataFinal[i]["_id"].(primitive.ObjectID)

		filter := bson.M{
			"email": bson.M{"$exists": false},
			//"nik":      dataFinal[i]["nik"],
			//"kk":       dataFinal[i]["kk"],
		}
		update := bson.M{
			"$set": bson.M{
				"role": 99,
			},
			//"$unset": bson.M{
			//	"fullName": "",
			//	"nik":      "",
			//	"kk":       "",
			//},
		}

		updated, err := db.Collection("users").UpdateMany(ctx, filter, update)
		if err != nil {
			return err
		}
		fmt.Println("data ke : ", i, updated)

	}
	return nil
}
