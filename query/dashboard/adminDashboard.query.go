package query

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"hutanku-service/config"
	"hutanku-service/models"
	"log"
	"time"
)

func GetHeaderDashboardData(c echo.Context) (models.Response, error) {
	var res models.Response
	var HeaderDashboardData models.HeaderDashboardAdmin
	chanJumlahPokja := make(chan []interface{})
	chanJumlahAnggota := make(chan int64)
	chanJumlahPetak := make(chan int64)
	chanJumlahAndil := make(chan int64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := config.Connect()
	if err != nil {
		return res, err
	}

	defer cancel()
	go func() {
		jumlahPokja, err := db.Collection("petak").Distinct(ctx, "pokja", bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		chanJumlahPokja <- jumlahPokja

		defer close(chanJumlahPokja)
	}()
	go func() {
		jumlahPetak, err := db.Collection("petak").Distinct(ctx, "petak", bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		chanJumlahPetak <- int64(len(jumlahPetak))

		defer close(chanJumlahPetak)
	}()
	go func() {
		jumlahAndil, err := db.Collection("petak").Distinct(ctx, "andil", bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		chanJumlahAndil <- int64(len(jumlahAndil))

		defer close(chanJumlahAndil)
	}()
	go func() {
		filterJumlahAnggota := bson.M{"role": 99}
		jumlahAnggota, err := db.Collection("users").CountDocuments(ctx, filterJumlahAnggota)
		if err != nil {
			log.Fatal(err)
		}
		chanJumlahAnggota <- jumlahAnggota

		defer close(chanJumlahAnggota)
	}()

	jumlahPokja := <-chanJumlahPokja

	var jumlahAnggotaPokjaTemp []bson.M
	for i := 0; i < len(jumlahPokja); i++ {
		jumlahAnggotaPokja, err := db.Collection("petak").Distinct(ctx, "userId", bson.M{"pokja": jumlahPokja[i]})
		if err != nil {
			log.Fatal(err)
		}

		matchStage := bson.D{
			{"$match", bson.D{
				{"pokja", jumlahPokja[i]},
			}},
		}
		groupStage := bson.D{
			{"$group", bson.D{
				{"_id", "$pokja"},
				{"luasLahan", bson.D{
					{"$sum", "$luasLahan"},
				}},
			}},
		}
		cursor, err := db.Collection("petak").Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
		if err != nil {
			log.Fatal(err)
		}
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}

		jumlahAnggotaPokjaTemp = append(jumlahAnggotaPokjaTemp, bson.M{
			jumlahPokja[i].(string): bson.M{
				"jumlahAnggota":    len(jumlahAnggotaPokja),
				"luasLahanGarapan": results[0]["luasLahan"],
			},
		})
	}

	HeaderDashboardData.JumlahPokja = int64(len(jumlahPokja))
	HeaderDashboardData.IsiPokja = jumlahPokja
	HeaderDashboardData.JumlahAnggota = <-chanJumlahAnggota
	HeaderDashboardData.JumlahPetak = <-chanJumlahPetak
	HeaderDashboardData.JumlahAndil = <-chanJumlahAndil
	HeaderDashboardData.JumlahAnggotaPerPokja = jumlahAnggotaPokjaTemp

	res.Message = "Get data success"
	res.Data = HeaderDashboardData

	return res, nil
}
