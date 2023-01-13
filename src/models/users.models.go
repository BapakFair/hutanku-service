package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName     string             `json:"fullName" bson:"fullName" validate:"required, max=50"`
	NomorAnggota string             `json:"nomorAnggota" bson:"nomorAnggota" validate:"required, max=50"`
	TotalPetak   int                `json:"totalPetak" bson:"totalPetak" validate:"required, max=2000"`
	TotalAndil   int                `json:"totalAndil" bson:"totalAndil" validate:"required, max=2000"`
	Email        string             `json:"email" bson:"email" validate:"required"`
	Password     string             `json:"password" bson:"password" validate:"required"`
	PhoneNumber  string             `json:"phoneNumber" bson:"phoneNumber" validate:"required"`
	Dusun        string             `json:"dusun" bson:"dusun"`
	Desa         string             `json:"desa" bson:"desa"`
	RT           int                `json:"rt" bson:"rt"`
	RW           int                `json:"rw" bson:"rw"`
	Kecamatan    string             `json:"kecamatan" bson:"kecamatan"`
	Kabupaten    string             `json:"kabupaten" bson:"kabupaten"`
	Kelurahan    string             `json:"kelurahan" bson:"kelurahan"`
	Kota         string             `json:"kota" bson:"kota"`
	Alamat       string             `json:"alamat" bson:"alamat"`
	NIK          string             `json:"nik" bson:"nik"`
	KK           string             `json:"kk" bson:"kk"`
	Province     string             `json:"province" bson:"province"`
	Pokja        string             `json:"pokja" bson:"pokja"`
	Role         int                `json:"role" bson:"role"`
	ResetToken   string             `json:"resetToken" bson:"resetToken"`
}

type QueryGetUsers struct {
	ID           string `json:"id" bson:"id"`
	NomorAnggota string `json:"nomorAnggota" bson:"nomorAnggota"`
}

type ResetPassword struct {
	Token              string `json:"token" bson:"token"`
	NewPassword        string `json:"newPassword" bson:"newPassword"`
	NewPasswordConfirm string `json:"newPasswordConfirm" bson:"newPasswordConfirm"`
}
