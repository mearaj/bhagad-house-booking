package service

import (
	"errors"
	"github.com/mearaj/bhagad-house-booking/common/alog"
	. "github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
	"sync"
)

type Service interface {
	Initialized() bool
	Booking() Booking
	Bookings() <-chan []Booking
	Customers(bookingPublicKey string, offset, limit int) <-chan []Customer
	CreateBooking(privateKeyHex string) <-chan error
	SaveCustomer(contactPublicKey string, identified bool) <-chan error
	AutoCreateBooking() <-chan error
	BookingKeyExists(publicKey string) <-chan bool
	SetAsCurrentBooking(booking Booking) <-chan error
	Subscribe(topics ...EventTopic) Subscriber
	DeleteBookings([]Booking) <-chan error
	DeleteCustomers([]Customer) <-chan error
	CustomersCount(addrPublicKey string) <-chan int64
	BookingsCount() <-chan int64
}

// Service Always call GetServiceInstance function to create Service
type service struct {
	booking       Booking
	bookingMutex  sync.RWMutex
	database      interface{}
	databaseMutex sync.RWMutex
	eventBroker   *eventBroker
}

var serviceInstance = service{
	eventBroker: newEventBroker(),
}

func init() {
	go serviceInstance.init()
}

func GetServiceInstance() Service {
	return &serviceInstance
}

func (s *service) Booking() Booking {
	s.bookingMutex.RLock()
	defer s.bookingMutex.RUnlock()
	return s.booking
}

func (s *service) setBooking(booking Booking) {
	currBooking := s.Booking()
	s.bookingMutex.Lock()
	s.booking = booking
	s.bookingMutex.Unlock()
	if currBooking.ID != booking.ID &&
		booking.ID != 0 {
		<-s.saveBookingToDB(booking)
		event := Event{Data: BookingChangedEventData{}, Topic: BookingChangedEventTopic}
		s.eventBroker.Fire(event)
	}
}

func (s *service) Initialized() bool {
	s.databaseMutex.RLock()
	defer s.databaseMutex.RUnlock()
	return s.database != nil
}

func (s *service) CreateBooking(pvtKeyHex string) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
		if strings.TrimSpace(pvtKeyHex) == "" {
			err = errors.New("private key is empty")
			return
		}
		if !s.Initialized() {
			err = errors.New("database engine not running")
			return
		}
		booking := Booking{}
		s.setBooking(booking)
	}()
	return errCh
}

func (s *service) AutoCreateBooking() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
		if !s.Initialized() {
			err = errors.New("database engine not running")
			return
		}
		booking := Booking{}
		s.setBooking(booking)
		s.eventBroker.Fire(Event{
			Data:  BookingsChangedEventData{},
			Topic: BookingsChangedEventTopic,
		})
	}()
	return errCh
}

func (s *service) Subscribe(topics ...EventTopic) Subscriber {
	subscr := newSubscriber()
	_ = subscr.Subscribe(topics...)
	s.eventBroker.addSubscriber(subscr)
	return subscr
}

func recoverPanic(entry *logrus.Entry) {
	if r := recover(); r != nil {
		entry.Errorln("recovered from panic", r)
	}
}
func recoverPanicCloseCh[S any](stateChan chan<- S, state S, entry *logrus.Entry) {
	recoverPanic(entry)
	stateChan <- state
	close(stateChan)
}

func (s *service) init() {
	err := <-s.openDatabase()
	if err != nil {
		log.Fatal(err)
	}
	bookings := <-s.reqBookingsFromDB()
	if len(bookings) > 0 {
		s.setBooking(bookings[0])
	}
}

// saveBookingToDB saves as primary booking
func (s *service) saveBookingToDB(acc Booking) <-chan error {
	errCh := make(chan error, 1)
	var err error
	go func() {
		defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
	}()
	return errCh
}

func (s *service) reqBookingsFromDB() <-chan []Booking {
	bookingsCh := make(chan []Booking, 1)
	go func() {
		bookings := make([]Booking, 0)
		defer func() {
			recoverPanicCloseCh(bookingsCh, bookings, alog.Logger())
		}()
	}()
	return bookingsCh
}

func (s *service) Bookings() <-chan []Booking {
	bookingsCh := make(chan []Booking, 1)
	bookings := make([]Booking, 0)
	go func() {
		defer func() {
			recoverPanicCloseCh(bookingsCh, bookings, alog.Logger())
		}()
	}()
	return bookingsCh
}

func (s *service) BookingKeyExists(publicKey string) <-chan bool {
	existsCh := make(chan bool, 1)
	go func() {
		exists := false
		defer func() {
			recoverPanicCloseCh(existsCh, exists, alog.Logger())
		}()
	}()
	return existsCh
}

func (s *service) Customers(bookingPublicKey string, offset, limit int) <-chan []Customer {
	contactsCh := make(chan []Customer, 1)
	contacts := make([]Customer, 0)
	go func() {
		defer func() {
			recoverPanicCloseCh(contactsCh, contacts, alog.Logger())
		}()
	}()
	return contactsCh
}

func (s *service) SaveCustomer(publicKey string, identified bool) <-chan error {
	errCh := make(chan error, 1)
	a := s.Booking()
	if a.ID == 0 {
		errCh <- errors.New("current id is nil")
		close(errCh)
		return errCh
	}
	go func() {
		var err error
		defer func() {
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
		eventData := CustomersChangeEventData{}
		event := Event{Data: eventData, Topic: CustomersChangedEventTopic}
		s.eventBroker.Fire(event)
	}()
	return errCh
}

func (s *service) SetAsCurrentBooking(booking Booking) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
		s.setBooking(booking)
	}()
	return errCh
}

func (s *service) DeleteBookings(bookings []Booking) <-chan error {
	errCh := make(chan error, 1)
	if len(bookings) == 0 {
		errCh <- errors.New("bookings is empty")
		close(errCh)
		return errCh
	}
	go func() {
		var err error
		defer func() {
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
		acc := s.Booking()
		var currentBookingDeleted bool
		for _, eachBooking := range bookings {
			if acc.ID == eachBooking.ID {
				currentBookingDeleted = true
			}

		}
		accs := <-s.reqBookingsFromDB()
		if currentBookingDeleted {
			if len(accs) > 0 {
				s.setBooking(accs[0])
			} else {
				s.setBooking(Booking{})
			}
			s.eventBroker.Fire(Event{
				Data:  CustomersChangeEventData{},
				Topic: CustomersChangedEventTopic,
			})
		}
		s.eventBroker.Fire(Event{
			Data:  BookingsChangedEventData{},
			Topic: BookingsChangedEventTopic,
		})
	}()
	return errCh
}

func (s *service) BookingsCount() <-chan int64 {
	countCh := make(chan int64, 1)
	go func() {
		count := int64(0)
		defer func() {
			recoverPanicCloseCh(countCh, count, alog.Logger())
		}()
	}()
	return countCh
}

func (s *service) CustomersCount(bookingPublicKey string) <-chan int64 {
	countCh := make(chan int64, 1)
	go func() {
		count := int64(0)
		defer func() {
			recoverPanicCloseCh(countCh, count, alog.Logger())
		}()
	}()
	return countCh
}

func (s *service) DeleteCustomers(contacts []Customer) <-chan error {
	errCh := make(chan error, 1)
	if len(contacts) == 0 {
		errCh <- errors.New("customers is empty")
		close(errCh)
		return errCh
	}
	go func() {
		var err error
		defer func() {
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
	}()
	return errCh
}
