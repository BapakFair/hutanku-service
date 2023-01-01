package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mailgun/mailgun-go/v4"
	"go.mongodb.org/mongo-driver/bson"
	"hutanku-service/config"
	helper "hutanku-service/helpers"
	"hutanku-service/models"
	"log"
	"os"
	"time"
)

func ForgotPassword(c echo.Context) (models.Response, error) {
	var res models.Response
	var reqBody models.Users
	if err := c.Bind(&reqBody); err != nil {
		log.Fatal(err.Error())
	}

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
	body := "Silahkan klik link berikut ini untuk mengubah password anda : " + domainBase + "/" + resetToken
	recipient := reqBody.Email
	message := mg.NewMessage(sender, subject, body, recipient)
	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	// =================================================

	res.Message = "Please Check Your Email"
	res.Data = result.ModifiedCount

	return res, nil
}

func ResetPassword(c echo.Context) (models.Response, error) {
	var res models.Response
	var users models.Users
	var reqBody models.ResetPassword
	if err := c.Bind(&reqBody); err != nil {
		log.Fatal(err.Error())
	}

	// set config mailgun
	privateApi := os.Getenv("API_KEY_MAILGUN")
	domainMail := os.Getenv("DOMAIN_MAIL")
	domainBase := os.Getenv("DOMAIN_BASE")
	mg := mailgun.NewMailgun(domainMail, privateApi)
	//====================

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	db, err := config.Connect()
	if err != nil {
		return res, err
	}
	defer cancel()

	if reqBody.NewPassword != reqBody.NewPasswordConfirm {
		err := errors.New("password dan confirmasi password tidak sama")
		return res, err
	}
	newPassword, _ := helper.HashPassword(reqBody.NewPassword)

	// get data user for getting their email
	if err := db.Collection("users").FindOne(ctx, bson.M{"resetToken": reqBody.Token}).Decode(&users); err != nil {
		return res, err
	}

	// set resetToken to empty string & update password on database user
	filter := bson.M{"resetToken": reqBody.Token}
	update := bson.M{
		"$set": bson.M{
			"resetToken": "",
			"password":   newPassword,
		},
	}
	result, err := db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		return res, err
	}
	// ============================

	// send email confirmation with token to user email
	sender := "superman@hutanku.com"
	subject := "Password anda berhasil diganti"
	body := "Berikut adalah password baru anda :" + newPassword + "silahkan login pada halaman berikut : " + domainBase + "/login"
	recipient := users.Email
	message := mg.NewMessage(sender, subject, body, recipient)
	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		return res, err
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	// =================================================

	res.Message = "Please Check Your Email"
	res.Data = result.ModifiedCount

	return res, nil
}
