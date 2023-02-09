package service

import (
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/common/request"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/frontend"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
)

type (
	Booking                      = model.Booking
	CreateBookingRequest         = request.CreateBooking
	UpdateBookingRequest         = request.UpdateBooking
	DeleteBookingRequest         = request.DeleteBooking
	BookingsRequest              = request.ListBookings
	TransactionsRequest          = request.GetTransactions
	AddUpdateTransactionRequest  = request.AddUpdateTransaction
	DeleteTransactionRequest     = request.DeleteTransaction
	BookingsResponse             = response.Bookings
	SearchBookingsResponse       = response.SearchBookings
	CreateBookingResponse        = response.CreateBooking
	DeleteBookingResponse        = response.DeleteBooking
	UpdateBookingResponse        = response.UpdateBooking
	TransactionsResponse         = response.GetTransactions
	UserResponse                 = response.LoginUser
	AddUpdateTransactionResponse = response.AddUpdateTransaction
	DeleteTransactionResponse    = response.DeleteTransaction
	NewBookingSMSResponse        = response.NewBookingSMS
	NewBookingEmailResponse      = response.NewBookingEmail
	NewTransactionSMSResponse    = response.NewTransactionSMS
	NewTransactionEmailResponse  = response.NewTransactionEmail
)

type Service interface {
	LogInUser(email string, password string)
	LogOutUser()
	GetBookings(request BookingsRequest)
	SearchBookings(query string)
	CreateBooking(booking CreateBookingRequest)
	UpdateBooking(booking UpdateBookingRequest)
	Subscribe(topics ...Topic) Subscriber
	DeleteBooking(number int)
	GetTransactions(request TransactionsRequest, eventID interface{})
	AddUpdateTransaction(request AddUpdateTransactionRequest, eventID interface{})
	DeleteTransaction(request DeleteTransactionRequest, eventID interface{})
	SendNewBookingSMS(bookingNumber int, eventID interface{})
	SendNewBookingEmail(bookingNumber int, eventID interface{})
	SendNewTransactionSMS(transactionNumber int, eventID interface{})
	SendNewTransactionEmail(transactionNumber int, eventID interface{})
}

type service struct {
	config      frontend.Config
	eventBroker *eventBroker
}

func NewService() Service {
	s := &service{eventBroker: newEventBroker()}
	s.config = frontend.LoadConfig()
	user.LoadSettings()
	s.eventBroker.Fire(Event{
		Data:   user.User(),
		Topic:  TopicLoggedInOut,
		Cached: false,
		ID:     nil,
	})
	return s
}

func (s *service) Subscribe(topics ...Topic) Subscriber {
	subscr := newSubscriber()
	_ = subscr.Subscribe(topics...)
	s.eventBroker.addSubscriber(subscr)
	return subscr
}
