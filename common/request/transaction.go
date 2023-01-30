package request

import (
	"github.com/mearaj/bhagad-house-booking/common/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddUpdateTransaction = model.Transaction

type GetTransactions struct {
	BookingID primitive.ObjectID `json:"booking_id" bson:"booking_id" binding:"required,booking_id"`
}
type DeleteTransaction struct {
	ID        primitive.ObjectID `json:"_id"        bson:"_id"`
	BookingID primitive.ObjectID `json:"booking_id"        bson:"booking_id"`
}
