package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Petak struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	Petak     string             `json:"petak" bson:"petak"`
	Andil     string             `json:"andil" bson:"andil"`
	Pokja     string             `json:"pokja" bson:"pokja"`
	LuasLahan float64            `json:"luasLahan" bson:"luasLahan"`
}
