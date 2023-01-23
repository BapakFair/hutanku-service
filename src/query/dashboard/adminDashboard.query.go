package query

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hutanku-service/config"
	"hutanku-service/src/models"
	"log"
	"sync"
	"time"
)

func GetHeaderDashboardData() (models.Response, error) {
	var res models.Response
	var HeaderDashboardData models.HeaderDashboardAdmin
	chanJumlahPokja := make(chan []interface{})
	chanJumlahAnggota := make(chan int64)
	chanJumlahPetak := make(chan int64)
	chanJumlahAndil := make(chan int64)
	chanLocation := make(chan []bson.M)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := config.Connect()
	if err != nil {
		return res, err
	}

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

	go func() {
		filterLocation := bson.M{}
		isiLocation, err := db.Collection("location").Find(ctx, filterLocation)
		if err != nil {
			log.Fatal(err)
		}
		var dataLocation []bson.M
		if err := isiLocation.All(ctx, &dataLocation); err != nil {
			return
		}
		chanLocation <- dataLocation

		defer close(chanLocation)
	}()

	jumlahPokja := <-chanJumlahPokja
	var jumlahAnggotaPokjaTemp []bson.M

	for _, pokja := range jumlahPokja {
		wg := new(sync.WaitGroup)
		wg.Add(4)
		jumlahAnggotaPokja, err := db.Collection("petak").Distinct(ctx, "userId", bson.M{"pokja": pokja})
		if err != nil {
			log.Fatal(err)
		}

		var results []bson.M
		go func() {
			defer wg.Done()
			matchStage := bson.D{
				{"$match", bson.D{
					{"pokja", pokja},
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
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}
		}()

		var jumlahPetakPokja []interface{}
		go func() {
			defer wg.Done()
			res, err := db.Collection("petak").Distinct(ctx, "petak", bson.M{"pokja": pokja})
			if err != nil {
				log.Fatal(err)
			}
			jumlahPetakPokja = res
		}()

		var jumlahAndilPokja []interface{}
		go func() {
			defer wg.Done()
			res, err := db.Collection("petak").Distinct(ctx, "andil", bson.M{"pokja": pokja})
			if err != nil {
				log.Fatal(err)
			}
			jumlahAndilPokja = res
		}()

		var jumlahKaryawanPokja int64
		go func() {
			defer wg.Done()
			res, err := db.Collection("users").CountDocuments(ctx, bson.M{"pokja": pokja, "role": 11})
			if err != nil {
				log.Fatal(err)
			}
			jumlahKaryawanPokja = res
		}()
		wg.Wait()

		jumlahAnggotaPokjaTemp = append(jumlahAnggotaPokjaTemp, bson.M{
			pokja.(string): bson.M{
				"jumlahAnggota":    len(jumlahAnggotaPokja),
				"luasLahanGarapan": results[0]["luasLahan"],
				"jumlahPetak":      len(jumlahPetakPokja),
				"jumlahAndil":      len(jumlahAndilPokja),
				"jumlahKaryawan":   &jumlahKaryawanPokja,
			},
		})

	}

	var isiLocationBaru []bson.M
	isiLocation := <-chanLocation
	for _, location := range isiLocation {
		wg := new(sync.WaitGroup)
		wg.Add(4)
		petak := location["id"]

		var jumlahAndilTemp int64
		go func() {
			defer wg.Done()
			jumlahAndil, err := db.Collection("petak").CountDocuments(ctx, bson.M{"petak": bson.M{
				"$regex": primitive.Regex{
					Pattern: petak.(string),
					Options: "i",
				},
			}})
			if err != nil {
				log.Fatal(err)
			}
			jumlahAndilTemp = jumlahAndil
		}()

		var namaPokja []bson.M
		var isiPokja string

		go func() {
			defer wg.Done()
			data, err := db.Collection("petak").Find(ctx, bson.M{"petak": bson.M{
				"$regex": primitive.Regex{
					Pattern: petak.(string),
					Options: "i",
				},
			}},
			)
			if err != nil {
				log.Fatal(err)
			}
			if err := data.All(ctx, &namaPokja); err != nil {
				return
			}
			if len(namaPokja) > 0 {
				isiPokja = namaPokja[0]["pokja"].(string)
			} else {
				isiPokja = ""
			}
		}()

		var results []bson.M
		go func() {
			defer wg.Done()
			matchStage := bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "petak", Value: bson.D{
						{
							Key: "$regex",
							Value: primitive.Regex{
								Pattern: petak.(string),
								Options: "i",
							},
						},
					}},
				}},
			}
			groupStage := bson.D{
				{"$group", bson.D{
					{"_id", "$petak"},
					{"luasLahan", bson.D{
						{"$sum", "$luasLahan"},
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
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}

			if len(results) == 0 {
				results = append(results, bson.M{
					"luasLahan": 0,
				},
				)

			}
		}()

		var jumlahAnggota []bson.M
		go func() {
			defer wg.Done()
			matchStage := bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "petak", Value: bson.D{
						{
							Key: "$regex",
							Value: primitive.Regex{
								Pattern: petak.(string),
								Options: "i",
							},
						},
					}},
				}},
			}
			groupStage := bson.D{
				{"$group", bson.D{
					{"_id", "$userId"},
				}},
			}
			cursor, err := db.Collection("petak").Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
			if err != nil {
				log.Fatal(err)
			}
			if err = cursor.All(context.TODO(), &jumlahAnggota); err != nil {
				panic(err)
			}
		}()
		wg.Wait()

		isiLocationBaru = append(isiLocationBaru, bson.M{
			petak.(string): bson.M{
				"pokja":         isiPokja,
				"jumlahAnggota": len(jumlahAnggota),
				"luasLahan":     results[0]["luasLahan"],
				"name":          location["name"],
				"id":            &petak,
				"jumlahAndil":   &jumlahAndilTemp,
				"path":          location["path"],
			},
		})

	}

	HeaderDashboardData.JumlahPokja = int64(len(jumlahPokja))
	HeaderDashboardData.IsiPokja = jumlahPokja
	HeaderDashboardData.JumlahAnggota = <-chanJumlahAnggota
	HeaderDashboardData.JumlahPetak = <-chanJumlahPetak
	HeaderDashboardData.JumlahAndil = <-chanJumlahAndil
	HeaderDashboardData.JumlahAnggotaPerPokja = jumlahAnggotaPokjaTemp
	HeaderDashboardData.Location = isiLocationBaru

	res.Message = "Get data success"
	res.Data = HeaderDashboardData

	return res, nil
}
