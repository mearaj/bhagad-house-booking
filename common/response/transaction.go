package response

import (
	"github.com/mearaj/bhagad-house-booking/common/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddUpdateTransaction struct {
	Transaction model.Transaction `json:"transaction,omitempty" bson:"transaction,omitempty"`
	Error       string            `json:"error,omitempty" bson:"transactions,omitempty"`
}

type GetTransactions struct {
	Transactions []model.Transaction `json:"transactions,omitempty" bson:"transactions,omitempty"`
	Error        string              `json:"error,omitempty" bson:"transactions,omitempty"`
}

type DeleteTransaction struct {
	ID        primitive.ObjectID `json:"_id,omitempty"        bson:"_id,omitempty"`
	BookingID primitive.ObjectID `json:"booking_id,omitempty"        bson:"booking_id,omitempty"`
	Error     string             `json:"error,omitempty"`
}
