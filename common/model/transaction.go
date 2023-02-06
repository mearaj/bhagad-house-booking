package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transaction struct {
	ID            primitive.ObjectID `json:"_id,omitempty"        bson:"_id,omitempty"`
	BookingNumber int                `json:"booking_number"        bson:"booking_number"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
	Amount        float64            `json:"amount" bson:"amount"`
	Details       string             `json:"details"    bson:"details"`
}
