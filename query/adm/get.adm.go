package query

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hutanku-service/config"
	"hutanku-service/models"
	"time"
)

// function to get administration area Indonesian Data
func GetAdm(c echo.Context) (models.Response, error) {
	var res models.Response
	ctx, cancel := context.WithTimeout(context.Background(), 40000*time.Second)
	db, err := config.Connect()
	if err != nil {
		return res, err
	}

	defer cancel()
	var data []bson.M

	filter := bson.M{}
	if q := c.QueryParam("q"); q != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"desakel": bson.M{
						"$regex": primitive.Regex{
							Pattern: q,
							Options: "i",
						},
					},
				}, {
					"kecamatan": bson.M{
						"$regex": primitive.Regex{
							Pattern: q,
							Options: "i",
						},
					},
				},
			},
		}
	}
	dataPick, err := db.Collection("adm_desakel").Find(ctx, filter)
	//dataDes, err := db.Collection("adm_desakel").Find(ctx, bson.M{
	//	"propinsi": bson.M{
	//		"$exists": false,
	//	},
	//})
	if err != nil {
		return res, err
	}
	if err := dataPick.All(ctx, &data); err != nil {
		return res, err
	}

	if len(data) == 0 {
		err = errors.New("no documents in result")
		return res, err
	}

	// this line only call to modify new data adm.
	//helper.GetAdmDetail(dataDesa, ctx)
	// ===========================================

	res.Message = "Get data success"
	res.Data = data

	return res, nil
}
