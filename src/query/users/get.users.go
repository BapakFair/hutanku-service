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
	"hutanku-service/src/helpers"
	"hutanku-service/src/models"
	"log"
	"os"
	"strconv"
	"time"
)

func GetUsers(c echo.Context) (models.ResponseWithPagination, error) {
	// get users by query string "_id", "pokja", "fullName", "nomorAnggota".

	var res models.ResponseWithPagination
	dataLoop := make(chan []bson.M)
	totalData := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := config.Connect()
	if err != nil {
		return res, err
	}

	filter := bson.M{}
	findOptions := options.Find()
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("perPage"))
	findOptions.SetSkip((int64(page) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))
	//findOptions.SetProjection(bson.M{"_id": 1, "pokja": 1})
	var dataFinal []bson.M

	go func() {
		defer close(dataLoop)
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
			log.Fatal(err)
		}
		if err := data.All(ctx, &dataFinal); err != nil {
			log.Fatal(err)
		}
		dataLoop <- dataFinal
	}()

	dataFromChan := <-dataLoop
	if len(dataFromChan) == 0 {
		err = errors.New("no documents in result")
		return res, err
	}

	go func() {
		defer close(totalData)
		totalDataDB, err := db.Collection("users").CountDocuments(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		totalDataCount := int(totalDataDB) / perPage
		if totalDataCount%1 != 0 {
			totalDataCount = totalDataCount + 1
		}
		totalData <- totalDataCount
	}()

	// this line of code below to manual hash nik & kk data from string to hashed vice versa =======================
	// don't forget to change context timeout 120 second per 1000 data
	//err = helper.EncryptNikKk(dataFinal, ctx)
	//err = helper.DecryptNikKk(dataFinal, ctx)
	//err = helper.UpdateIdUserToPetak(dataFinal, ctx)
	//hasil, err := helper.ChangeDesaToDesakel(dataFinal, ctx)
	//if err != nil {
	//	return res, err
	//}
	//fmt.Println(hasil)
	// ==================================================================================================

	for _, data := range dataFromChan {
		secret := os.Getenv("RAHASIA_NEGARA")
		NIK := fmt.Sprintf("%v", data["nik"])
		KK := fmt.Sprintf("%v", data["kk"])
		dataKK := make(chan string)
		dataNIK := make(chan string)
		jumlahPetak := make(chan int64)
		lahanGarapan := make(chan []bson.M)

		Phone := fmt.Sprintf("%v", data["phoneNumber"])
		if Phone == "<nil>" {
			data["phoneNumber"] = ""
		} else {
			data["phoneNumber"] = string(helper.Decrypt([]byte(Phone), secret))
		}

		Alamat := fmt.Sprintf("%v", data["alamat"])
		fmt.Println(Alamat)
		if len(Alamat) < 8 {
			data["alamat"] = Alamat
		} else {
			data["alamat"] = string(helper.Decrypt([]byte(Alamat), secret))
		}
		go func() {
			defer close(dataNIK)
			dataNIK <- string(helper.Decrypt([]byte(NIK), secret))
		}()
		go func() {
			defer close(dataKK)
			dataKK <- string(helper.Decrypt([]byte(KK), secret))
		}()

		data["nik"] = <-dataNIK
		data["kk"] = <-dataKK

		checkPetak, err := db.Collection("petak").Find(ctx, bson.M{"userId": data["_id"]})
		if err != nil {
			return res, err
		}

		var dataPetak []bson.M
		if err := checkPetak.All(ctx, &dataPetak); err != nil {
			err := errors.New("binding cursor gagal")
			return res, err
		}

		if len(dataPetak) == 0 {
			data["jumlahPetak"] = 0
			data["totalLahanGarapan"] = 0
		} else {
			// menghitung jumlah petak yang dimiliki anggota =============
			go func() {
				defer close(jumlahPetak)
				jumlahPetakUser, err := db.Collection("petak").CountDocuments(ctx, bson.M{
					"pokja":  data["pokja"],
					"userId": data["_id"],
				})
				if err != nil {
					log.Fatal(err)
				}
				jumlahPetak <- jumlahPetakUser
			}()

			// ============================================================

			// menghitung total luas lahan garapan ========================
			go func() {
				defer close(lahanGarapan)
				matchStage := bson.D{
					{"$match", bson.D{
						{"userId", data["_id"]},
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
					log.Fatal(err)
				}
				var results []bson.M
				if err = cursor.All(context.TODO(), &results); err != nil {
					log.Fatal(err)
				}
				lahanGarapan <- results
			}()

			// =============================================================
			lahanGarapanDariChan := <-lahanGarapan
			data["jumlahPetak"] = <-jumlahPetak
			data["totalLahanGarapan"] = lahanGarapanDariChan[0]["luasLahan"]
		}
	}

	res.Message = "Get data success"
	res.Data = dataFinal
	res.Page = page
	res.TotalData = <-totalData

	return res, nil
}
