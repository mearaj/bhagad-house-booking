package request

import (
	"github.com/mearaj/bhagad-house-booking/common/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddUpdateTransaction = model.Transaction

type GetTransactions struct {
	BookingNumber int `json:"booking_number" bson:"booking_number" binding:"required,booking_id"`
}
type DeleteTransaction struct {
	ID            primitive.ObjectID `json:"_id"        bson:"_id"`
	BookingNumber int                `json:"booking_number" bson:"booking_number"`
}
