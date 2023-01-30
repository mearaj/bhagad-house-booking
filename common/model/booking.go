package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	ID           primitive.ObjectID `json:"_id,omitempty"        bson:"_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	StartDate    time.Time          `json:"start_date" bson:"start_date"`
	EndDate      time.Time          `json:"end_date"   bson:"end_date"`
	Details      string             `json:"details"    bson:"details"`
	CustomerName string             `json:"customer_name" bson:"customer_name"`
	RatePerDay   float64            `json:"rate_per_day" bson:"rate_per_day"`
}
