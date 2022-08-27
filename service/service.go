package service

import (
	"errors"
	"gioui.org/app"
	"github.com/mearaj/bhagad-house-booking/alog"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
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

func (s *service) GormDB() *gorm.DB {
	s.databaseMutex.RLock()
	defer s.databaseMutex.RUnlock()
	if s.database == (*gorm.DB)(nil) {
		return nil
	}
	if gormDB, ok := s.database.(*gorm.DB); ok {
		return gormDB
	}
	return nil
}

func (s *service) setGormDB(gormDB *gorm.DB) {
	s.databaseMutex.Lock()
	defer s.databaseMutex.Unlock()
	s.database = gormDB
}

func (s *service) openDatabase() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			if err != nil {
				alog.Logger().Errorln(err)
			}
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
		dirPath, err := app.DataDir()
		if err != nil {
			return
		}
		dirPath = filepath.Join(dirPath, DBPathCfgDir)
		if _, err = os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, 0700)
			if err != nil {
				return
			}
		}
		dbFullName := filepath.Join(dirPath, DBPathFileName)
		if _, err = os.Stat(dbFullName); os.IsNotExist(err) {
			var file *os.File
			file, err = os.OpenFile(
				dbFullName,
				os.O_CREATE|os.O_APPEND|os.O_RDWR,
				0700,
			)
			if err != nil {
				return
			}
			_ = file.Close()
		}
		gormDB, err := gorm.Open(sqlite.Open(dbFullName), &gorm.Config{})
		if err != nil {
			return
		}
		s.setGormDB(gormDB)
		err = s.GormDB().AutoMigrate(&Booking{}, &Customer{}, &KeyValue{})
		if err != nil {
			return
		}
	}()
	return errCh
}

// saveBookingToDB saves as primary booking
func (s *service) saveBookingToDB(acc Booking) <-chan error {
	errCh := make(chan error, 1)
	var err error
	go func() {
		defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
		txn := s.GormDB().Save(&acc)
		if txn.Error != nil {
			err = txn.Error
			alog.Logger().Errorln(err)
		}
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
		if s.GormDB() == nil {
			return
		}
		result := s.GormDB().Model(&Booking{}).Order("updated_at desc").Find(&bookings)
		if result.Error != nil {
			alog.Logger().Errorln(result.Error)
		}
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
		txn := s.GormDB().Order("updated_at desc").Find(&bookings)
		if txn.Error != nil {
			alog.Logger().Errorln(txn.Error)
		}
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
		txn := s.GormDB().Find(&Booking{}, "id = ?", publicKey)
		if txn.Error != nil {
			alog.Logger().Errorln(txn.Error)
		}
		exists = txn.Error == nil && txn.RowsAffected == 1
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
		txn := s.GormDB().Order("updated_at desc").Offset(offset).Limit(limit).Find(&contacts, "booking_public_key = ?", bookingPublicKey)
		if txn.Error != nil {
			alog.Logger().Errorln(txn.Error)
		}
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
		var txn *gorm.DB
		ct := Customer{Name: publicKey}
		txn = s.GormDB().Save(&ct)
		if txn.Error != nil {
			err = txn.Error
			alog.Logger().Errorln(txn)
			return
		}
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
		txn := s.GormDB().Find(&booking, "id = ?", booking.ID)
		if txn.Error != nil {
			err = txn.Error
		}
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
			s.GormDB().Where("id = ?", eachBooking.ID).Delete(&eachBooking)
			s.GormDB().Where("customer_id = ?", eachBooking.ID).Delete(&Customer{})
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
		if s.GormDB() == nil {
			return
		}
		txn := s.GormDB().Model(&Booking{}).Count(&count)
		if txn.Error != nil {
			alog.Logger().Println(txn.Error)
		}
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
		if s.GormDB() == nil {
			return
		}
		txn := s.GormDB().Model(&Customer{}).Where(map[string]interface{}{
			"booking_public_key": bookingPublicKey,
		}).Count(&count)
		if txn.Error != nil {
			alog.Logger().Println(txn.Error)
		}
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
		for _, eachCustomer := range contacts {
			s.GormDB().Where("booking_public_key = ? and public_key = ?", eachCustomer.Name, eachCustomer.Name).Delete(&Customer{})
		}
		s.eventBroker.Fire(Event{
			Data:  CustomersChangeEventData{},
			Topic: CustomersChangedEventTopic,
		})
	}()
	return errCh
}
