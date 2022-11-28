package query

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"go.mongodb.org/mongo-driver/bson"
	"hutanku-service/config"
	helper "hutanku-service/helpers"
	"hutanku-service/models"
	"log"
	"net/http"
	"os"
	"time"
)

func ForgotPassword(reqBody models.Users) (models.Response, error) {
	var res models.Response

	// set config mailgun
	privateApi := os.Getenv("API_KEY_MAILGUN")
	domainMail := os.Getenv("DOMAIN_MAIL")
	domainBase := os.Getenv("DOMAIN_BASE")
	mg := mailgun.NewMailgun(domainMail, privateApi)
	//====================
	fmt.Println("test key : ", privateApi)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	db, err := config.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cancel()

	resetToken, _ := helper.HashPassword(reqBody.Email + "secretTokenForResetPassword")

	// set token on database user
	filter := bson.M{"email": reqBody.Email}
	update := bson.M{
		"$set": bson.M{
			"resetToken": resetToken,
		},
	}
	result, err := db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err.Error())
	}
	// ============================

	// send email confirmation with token to user email
	sender := "superman@hutanku.com"
	subject := "Reset Password"
	body := "Silahkan klik link berikut ini untuk mengubah password anda :" + domainBase + "/" + resetToken
	recipient := "fairsulaiman@gmail.com" // please change this dynamic in the future
	message := mg.NewMessage(sender, subject, body, recipient)
	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	// =================================================

	res.Status = http.StatusOK
	res.Message = "Please Check Your Email"
	res.Data = result.ModifiedCount

	return res, nil
}
