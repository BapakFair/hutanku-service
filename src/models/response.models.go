package models

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseWithPagination struct {
	Message   string      `json:"message" bson:"message"`
	Data      interface{} `json:"data" bson:"data"`
	Page      int         `json:"page" bson:"page"`
	TotalData int         `json:"totalData" bson:"totalData"`
}

type ResError struct {
	Status int   `json:"status"`
	Data   error `json:"data"`
}
