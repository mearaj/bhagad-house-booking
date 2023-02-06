package request

import (
	"github.com/mearaj/bhagad-house-booking/common/model"
	"time"
)

type CreateBooking = model.Booking

type ListBookings struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type UpdateBooking = model.Booking

type DeleteBooking struct {
	Number int `json:"number"`
}
