package models

import "go.mongodb.org/mongo-driver/bson"

type HeaderDashboardAdmin struct {
	JumlahPokja           int64         `json:"jumlahPokja" bson:"jumlahPokja" validate:"required, max=50"`
	IsiPokja              []interface{} `json:"isiPokja" bson:"isiPokja" validate:"required, max=50"`
	JumlahKaryawan        int64         `json:"jumlahKaryawan" bson:"jumlahKaryawan" validate:"required, max=100"`
	JumlahAnggota         int64         `json:"jumlahAnggota" bson:"jumlahAnggota" validate:"required, max=50"`
	JumlahPetak           int64         `json:"jumlahPetak" bson:"jumlahPetak"`
	JumlahAndil           int64         `json:"jumlahAndil" bson:"jumlahAndil"`
	JumlahAnggotaPerPokja []bson.M      `json:"jumlahAnggotaPerPokja" bson:"jumlahAnggotaPerPokja"`
	Location              []bson.M      `json:"location" bson:"location"`
}
