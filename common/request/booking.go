package request

import (
	"github.com/mearaj/bhagad-house-booking/common/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CreateBooking = model.Booking

type ListBookings struct {
	StartDate time.Time `json:"start_date" bson:"start_date"`
	EndDate   time.Time `json:"end_date" bson:"end_date"`
}

type UpdateBooking struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	StartDate    time.Time          `json:"start_date" bson:"start_date"`
	EndDate      time.Time          `json:"end_date" bson:"end_date"`
	Details      string             `json:"details" bson:"details"`
	CustomerName string             `json:"customer_name" bson:"customer_name"`
	RatePerDay   float64            `json:"rate_per_day" bson:"rate_per_day"`
}

type DeleteBooking struct {
	ID primitive.ObjectID `json:"_id"        bson:"_id"`
}
