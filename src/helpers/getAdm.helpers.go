package helper

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"hutanku-service/config"
)

// this function is to set manual combine data kecamatan -> desa
func GetAdmDetail(data []bson.M, ctx context.Context) ([]bson.M, error) {
	db, err := config.Connect()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(data); i++ {
		var dataKecamatan []bson.M
		dataKec, err := db.Collection("adm_propinsi").Find(ctx, bson.M{
			"kode": data[i]["kode_propinsi"],
		})
		if err != nil {
			return nil, err
		}
		if err := dataKec.All(ctx, &dataKecamatan); err != nil {
			return nil, err
		}

		filterKecamatan := bson.M{
			"kode_propinsi": data[i]["kode_propinsi"],
		}
		update := bson.M{
			"$set": bson.M{
				"propinsi": dataKecamatan[0]["nama"],
			},
		}

		_, err = db.Collection("adm_desakel").UpdateMany(ctx, filterKecamatan, update)
		if err != nil {
			return nil, err
		}

		fmt.Println("data ke : ", i, "ini isi data propinsi : ", dataKecamatan)

	}
	return nil, nil
}

func ChangeDesaToDesakel(data []bson.M, ctx context.Context) ([]bson.M, error) {
	db, err := config.Connect()
	if err != nil {
		return nil, err
	}

	//for i, data := range data {
	filterKecamatan := bson.M{}
	update := bson.M{
		"$rename": bson.M{
			"desa": "desakel",
		},
	}

	_, err = db.Collection("users").UpdateMany(ctx, filterKecamatan, update)
	if err != nil {
		return nil, err
	}

	//	fmt.Println("data ke : ", i, "isinya ini data : ", data)
	//
	//}
	return nil, nil
}
