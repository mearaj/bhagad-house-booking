package response

import (
	"github.com/mearaj/bhagad-house-booking/common/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddUpdateTransaction struct {
	Transaction model.Transaction `json:"transaction,omitempty"`
	Error       string            `json:"error,omitempty"`
}

type GetTransactions struct {
	Transactions []model.Transaction `json:"transactions,omitempty"`
	Error        string              `json:"error,omitempty"`
}

type DeleteTransaction struct {
	ID            primitive.ObjectID `json:"_id,omitempty"`
	BookingNumber int                `json:"booking_number,omitempty"`
	Error         string             `json:"error,omitempty"`
}

type NewTransactionEmail struct {
	Transaction model.Transaction `json:"transaction,omitempty"`
	Error       string            `json:"error,omitempty"`
}
type NewTransactionSMS struct {
	Transaction model.Transaction `json:"transaction,omitempty"`
	Error       string            `json:"error,omitempty"`
}
