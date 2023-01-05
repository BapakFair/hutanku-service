package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Propinsi struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Kode int                `json:"kode" bson:"kode" validate:"required, max=50"`
	Nama string             `json:"nama" bson:"nama" validate:"required, max=100"`
}

type KotaKab struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Kode         int                `json:"kode" bson:"kode" validate:"required, max=50"`
	KodePropinsi int                `json:"kode_propinsi" bson:"kode_propinsi" validate:"required, max=50"`
	Nama         string             `json:"nama" bson:"nama" validate:"required, max=100"`
}

type Kecamatan struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Kode          int                `json:"kode" bson:"kode" validate:"required, max=50"`
	KodeKecamatan int                `json:"kode_kecamatan" bson:"kode_kecamatan" validate:"required, max=50"`
	Nama          string             `json:"nama" bson:"nama" validate:"required, max=100"`
}

type DesaKel struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Kode          int                `json:"kode" bson:"kode" validate:"required, max=50"`
	KodeDesaKel   int                `json:"kode_desakel" bson:"kode_desakel" validate:"required, max=50"`
	Desakel       string             `json:"desakel" bson:"desakel" validate:"required, max=100"`
	KodeKecamatan int                `json:"kode_kecamatan" bson:"kode_kecamatan" validate:"required, max=50"`
	Kecamatan     string             `json:"kecamatan" bson:"kecamatan" validate:"required, max=100"`
	KodeKotaKab   int                `json:"kode_kotakab" bson:"kode_kotakab" validate:"required, max=50"`
	KotaKab       string             `json:"kotakab" bson:"kotakab" validate:"required, max=100"`
	KodePropinsi  int                `json:"kode_propinsi" bson:"kode_propinsi" validate:"required, max=50"`
	Propinsi      string             `json:"propinsi" bson:"propinsi" validate:"required, max=100"`
}
