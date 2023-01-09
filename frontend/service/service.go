package service

import (
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/frontend"
)

type (
	Booking                = sqlc.Booking
	BookingParams          = sqlc.ListBookingsParams
	BookingsResponse       = sqlc.BookingsResponse
	SearchBookingsResponse = sqlc.SearchBookingsResponse
	CreateBookingResponse  = sqlc.CreateBookingResponse
	DeleteBookingResponse  = sqlc.DeleteBookingResponse
	UpdateBookingResponse  = sqlc.UpdateBookingResponse
)

type Service interface {
	LogInUser(email string, password string)
	LogOutUser()
	Bookings(bookingParams BookingParams)
	SearchBookings(query string)
	CreateBooking(booking sqlc.CreateBookingParams)
	UpdateBooking(booking sqlc.UpdateBookingParams)
	Subscribe(topics ...Topic) Subscriber
	DeleteBooking(bookingID int64)
}

type service struct {
	config      frontend.Config
	eventBroker *eventBroker
}

func NewService() Service {
	s := &service{eventBroker: newEventBroker()}
	s.config = frontend.LoadConfig()
	return s
}

func (s *service) Subscribe(topics ...Topic) Subscriber {
	subscr := newSubscriber()
	_ = subscr.Subscribe(topics...)
	s.eventBroker.addSubscriber(subscr)
	return subscr
}
