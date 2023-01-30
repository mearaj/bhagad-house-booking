package service

import (
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/common/request"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/frontend"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
)

type Service interface {
	LogInUser(email string, password string)
	LogOutUser()
	Bookings(request BookingsRequest)
	SearchBookings(query string)
	CreateBooking(booking CreateBookingRequest)
	UpdateBooking(booking UpdateBookingRequest)
	Subscribe(topics ...Topic) Subscriber
	DeleteBooking(bookingID primitive.ObjectID)
	Transactions(request TransactionsRequest, eventID interface{})
	AddUpdateTransaction(request AddUpdateTransactionRequest, eventID interface{})
	DeleteTransaction(request DeleteTransactionRequest, eventID interface{})
}

type service struct {
	config       frontend.Config
	eventBroker  *eventBroker
	subscription *subscriber
}

func NewService() Service {
	s := &service{eventBroker: newEventBroker()}
	s.config = frontend.LoadConfig()
	s.subscription = newSubscriber()
	_ = s.subscription.Subscribe(TopicLoggedInOut)
	s.subscription.SubscribeWithCallback(s.onStateChange)
	s.eventBroker.addSubscriber(s.subscription)
	user.LoadSettings()
	s.eventBroker.Fire(Event{
		Data:   *user.User(),
		Topic:  TopicLoggedInOut,
		Cached: false,
		ID:     nil,
	})
	return s
}
func (s *service) onStateChange(event Event) {
	switch eventData := event.Data.(type) {
	case UserResponse:
		*user.User() = eventData
		user.SaveSettings()
	}
}

func (s *service) Subscribe(topics ...Topic) Subscriber {
	subscr := newSubscriber()
	_ = subscr.Subscribe(topics...)
	s.eventBroker.addSubscriber(subscr)
	return subscr
}
