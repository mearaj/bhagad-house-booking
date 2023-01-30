package response

import (
	"github.com/mearaj/bhagad-house-booking/common/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bookings struct {
	Bookings []model.Booking `json:"bookings,omitempty"`
	Error    string          `json:"error,omitempty"`
}

type AuthError struct {
	Error string `json:"error,omitempty"`
}

type CreateBooking struct {
	Booking model.Booking `json:"booking,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type UpdateBooking struct {
	Booking model.Booking `json:"booking,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type DeleteBooking struct {
	ID    primitive.ObjectID `json:"_id,omitempty"        bson:"_id,omitempty"`
	Error string             `json:"error,omitempty"`
}

type SearchBookings struct {
	Bookings []model.Booking `json:"bookings,omitempty"`
	Error    string          `json:"error,omitempty"`
}
