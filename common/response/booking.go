package response

import (
	"github.com/mearaj/bhagad-house-booking/common/model"
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

type DeleteBooking struct {
	Number int    `json:"number,omitempty"`
	Error  string `json:"error,omitempty"`
}
type UpdateBooking struct {
	Booking model.Booking `json:"booking,omitempty"`
	Error   string        `json:"error,omitempty"`
}
type SearchBookings struct {
	Bookings []model.Booking `json:"bookings,omitempty"`
	Error    string          `json:"error,omitempty"`
}

type NewBookingEmail struct {
	Booking model.Booking `json:"booking,omitempty"`
	Error   string        `json:"error,omitempty"`
}
type NewBookingSMS struct {
	Booking model.Booking `json:"booking,omitempty"`
	Error   string        `json:"error,omitempty"`
}
