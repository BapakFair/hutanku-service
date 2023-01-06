package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hutanku-service/config"
	helper "hutanku-service/helpers"
	"hutanku-service/models"
	"os"
	"strconv"
	"time"
)

// get users by query string "_id", "pokja", "fullName", "nomorAnggota".
func GetUsers(c echo.Context) (models.ResponseWithPagination, error) {
	var res models.ResponseWithPagination
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
	//findOptions.SetProjection(bson.M{"_id": 1, "pokja": 1})

	if q := c.QueryParam("q"); q != "" {
		objId, _ := primitive.ObjectIDFromHex(q)

		filter = bson.M{
			"role": 99,
			"$or": []bson.M{
				{
					"_id": objId,
				}, {
					"nomorAnggota": bson.M{
						"$regex": primitive.Regex{
							Pattern: q,
							Options: "i",
						},
					},
				}, {
					"fullName": bson.M{
						"$regex": primitive.Regex{
							Pattern: q,
							Options: "i",
						},
					},
				}, {
					"pokja": bson.M{
						"$regex": primitive.Regex{
							Pattern: q,
							Options: "i",
						},
					},
				},
			},
		}
	} else if q == "" {
		filter = bson.M{
			"role": 99,
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

	if len(dataFinal) == 0 {
		err = errors.New("no documents in result")
		return res, err
	}

	totalData, err := db.Collection("users").CountDocuments(ctx, filter)
	if err != nil {
		return res, err
	}
	totalDataCount := int(totalData) / perPage
	if totalDataCount%1 != 0 {
		totalDataCount = totalDataCount + 1
	}

	// this line of code below to manual hash nik & kk data from string to hashed vice versa =======================
	// don't forget to change context timeout 120 second per 1000 data
	//err = helper.EncryptNikKk(dataFinal, ctx)
	//err = helper.DecryptNikKk(dataFinal, ctx)
	//err = helper.UpdateIdUserToPetak(dataFinal, ctx)
	//if err != nil {
	//	return models.Response{}, err
	//}
	// ==================================================================================================

	for i := 0; i < len(dataFinal); i++ {

		secret := os.Getenv("RAHASIA_NEGARA")
		NIK := fmt.Sprintf("%v", dataFinal[i]["nik"])
		KK := fmt.Sprintf("%v", dataFinal[i]["kk"])
		//Phone := fmt.Sprintf("%v", dataFinal[i]["phoneNumber"])
		//Alamat := fmt.Sprintf("%v", dataFinal[i]["alamat"])

		// menghitung jumlah petak yang dimiliki anggota =============
		jumlahPetakUser, err := db.Collection("petak").CountDocuments(ctx, bson.M{
			"pokja":  dataFinal[i]["pokja"],
			"userId": dataFinal[i]["_id"],
		})
		if err != nil {
			return res, err
		}
		// ============================================================

		// menghitung total luas lahan garapan ========================
		matchStage := bson.D{
			{"$match", bson.D{
				{"userId", dataFinal[i]["_id"]},
				{"pokja", bson.D{{"$exists", true}}},
			}},
		}
		groupStage := bson.D{
			{"$group", bson.D{
				{"_id", "$pokja"},
				{"luasLahan", bson.D{
					{"$sum", bson.D{
						{"$toDecimal", "$luasLahan"},
					}},
				}},
			}},
		}

		projectStage := bson.D{
			{"$project", bson.D{
				{"luasLahan", bson.D{
					{"$round", bson.A{
						"$luasLahan", 2,
					}},
				}},
			}},
		}
		cursor, err := db.Collection("petak").Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		if err != nil {
			return res, err
		}
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			return res, err
		}

		// =============================================================
		dataFinal[i]["nik"] = string(helper.Decrypt([]byte(NIK), secret))
		dataFinal[i]["kk"] = string(helper.Decrypt([]byte(KK), secret))
		dataFinal[i]["jumlahPetak"] = jumlahPetakUser
		dataFinal[i]["totalLahanGarapan"] = results[0]["luasLahan"]
		//fmt.Println(i, " ini tersangka usernya: ", dataFinal[i], "-> ", results[0]["luasLahan"])
		//dataFinal[i]["phoneNumber"] = string(helper.Decrypt([]byte(Phone), secret))
		//dataFinal[i]["alamat"] = string(helper.Decrypt([]byte(Alamat), secret))
	}
	res.Message = "Get data success"
	res.Data = dataFinal
	res.Page = page
	res.TotalData = int(totalData)

	return res, nil
}
