package query

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"hutanku-service/config"
	"hutanku-service/src/helpers"
	models2 "hutanku-service/src/models"
	"log"
	"os"
	"time"
)

func CheckLogin(email, password string) (models2.Response, error) {
	var user models2.Users
	var res models2.Response

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	db, err := config.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cancel()

	if err := db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		err = errors.New("email or password keliru bosku ...")
		return res, err
	}

	match, err := helper.CheckPasswordHash(password, user.Password)
	if !match {
		err = errors.New("email or password keliru bosku ...")
		return res, err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["id"] = user.ID
	claims["nomorAnggota"] = user.NomorAnggota
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	kunci := os.Getenv("KUNCI_MASUK")

	t, err := token.SignedString([]byte(kunci))
	if err != nil {
		return res, err
	}

	res.Message = "login success"
	res.Data = t

	return res, nil
}
