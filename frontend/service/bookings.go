package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mearaj/bhagad-house-booking/common/alog"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/common/token"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
)

const UnauthorizedErrorStr = "unauthorized"

func (s *service) Bookings(bookingParams BookingsRequest) {
	go func() {
		var bookingsResponse BookingsResponse
		bookingsResponse.Bookings = make([]Booking, 0)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				bookingsResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  bookingsResponse,
				Topic: TopicFetchBookings,
			})
		}()
		params := url.Values{}
		startDate := bookingParams.StartDate.Format("2006-01-02")
		endDate := bookingParams.EndDate.Format("2006-01-02")
		params.Add("start_date", startDate)
		params.Add("end_date", endDate)
		req, err := http.NewRequest("GET", s.config.ApiURL+"/bookings?"+params.Encode(), nil)
		if err != nil {
			bookingsResponse.Error = err.Error()
			return
		}
		req.Header.Add("Accept", "application/json")
		userEvent, ok := s.eventBroker.cachedEvents.Get(TopicLoggedInOut)
		if ok {
			userResponse, ok := userEvent.Data.(UserResponse)
			if ok && userResponse.AccessToken != "" {
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userResponse.AccessToken))
			}
		}
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if err != nil {
			if errors.Is(err, token.ErrExpiredToken) || errors.Is(err, token.ErrInvalidToken) {
				s.eventBroker.Fire(Event{Data: UserResponse{}, Topic: TopicLoggedInOut})
			}
			bookingsResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&bookingsResponse)
		if err != nil {
			bookingsResponse.Error = err.Error()
			return
		}
	}()
}
func (s *service) SearchBookings(query string) {
	go func() {
		var searchBookingsResponse SearchBookingsResponse
		searchBookingsResponse.Bookings = make([]Booking, 0)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				searchBookingsResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  searchBookingsResponse,
				Topic: TopicSearchBookings,
			})
		}()
		params := url.Values{}
		params.Add("query", query)
		req, err := http.NewRequest("GET", s.config.ApiURL+"/bookings/search?"+params.Encode(), nil)
		if err != nil {
			searchBookingsResponse.Error = err.Error()
			return
		}
		req.Header.Add("Accept", "application/json")
		userEvent, ok := s.eventBroker.cachedEvents.Get(TopicLoggedInOut)
		if ok {
			userResponse, ok := userEvent.Data.(UserResponse)
			if ok && userResponse.AccessToken != "" {
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userResponse.AccessToken))
			}
		}
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if err != nil {
			if errors.Is(err, token.ErrExpiredToken) || errors.Is(err, token.ErrInvalidToken) {
				s.eventBroker.Fire(Event{Data: UserResponse{}, Topic: TopicLoggedInOut})
			}
			searchBookingsResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&searchBookingsResponse)
		if err != nil {
			searchBookingsResponse.Error = err.Error()
			return
		}
	}()
}
func (s *service) CreateBooking(bookingParams CreateBookingRequest) {
	go func() {
		var createBookingResponse CreateBookingResponse
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				createBookingResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  createBookingResponse,
				Topic: TopicCreateBooking,
			})
		}()
		userEvent, ok := s.eventBroker.cachedEvents.Get(TopicLoggedInOut)
		if !ok {
			createBookingResponse.Error = "user not logged in"
			return
		}
		userResponse, ok := userEvent.Data.(UserResponse)
		if !ok {
			createBookingResponse.Error = "critical error, need to contact admin"
			return
		}
		if userResponse.AccessToken == "" {
			createBookingResponse.Error = "user not logged in"
			return
		}
		jsonValues, err := json.Marshal(bookingParams)
		if err != nil {
			createBookingResponse.Error = err.Error()
			return
		}
		req, err := http.NewRequest("POST", s.config.ApiURL+"/bookings", bytes.NewBuffer(jsonValues))
		if err != nil {
			createBookingResponse.Error = err.Error()
			return
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userResponse.AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if rsp.StatusCode == http.StatusUnauthorized {
			userResponse = UserResponse{AccessToken: "", User: response.User{}, Error: ""}
			s.eventBroker.Fire(Event{Data: userResponse,
				Topic: TopicLoggedInOut,
			})
			createBookingResponse.Error = UnauthorizedErrorStr
			return
		}
		if err != nil {
			if errors.Is(err, token.ErrExpiredToken) || errors.Is(err, token.ErrInvalidToken) {
				s.eventBroker.Fire(Event{Data: UserResponse{}, Topic: TopicLoggedInOut})
			}
			createBookingResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&createBookingResponse)
		if err != nil {
			createBookingResponse.Error = err.Error()
			return
		}
	}()
}

func (s *service) UpdateBooking(bookingParams UpdateBookingRequest) {
	go func() {
		var updateBookingsResponse UpdateBookingResponse
		updateBookingsResponse.Booking.ID = bookingParams.ID
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				updateBookingsResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  updateBookingsResponse,
				Topic: TopicUpdateBooking,
			})
		}()
		userEvent, ok := s.eventBroker.cachedEvents.Get(TopicLoggedInOut)
		if !ok {
			updateBookingsResponse.Error = "user not logged in"
			return
		}
		userResponse, ok := userEvent.Data.(UserResponse)
		if !ok {
			updateBookingsResponse.Error = "critical error, need to contact admin"
			return
		}
		if userResponse.AccessToken == "" {
			updateBookingsResponse.Error = "user not logged in"
			return
		}
		jsonValues, err := json.Marshal(bookingParams)
		if err != nil {
			updateBookingsResponse.Error = err.Error()
			return
		}
		req, err := http.NewRequest("PUT", s.config.ApiURL+"/bookings", bytes.NewBuffer(jsonValues))
		if err != nil {
			updateBookingsResponse.Error = err.Error()
			return
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userResponse.AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if rsp.StatusCode == http.StatusUnauthorized {
			userResponse = UserResponse{AccessToken: "", User: response.User{}, Error: ""}
			s.eventBroker.Fire(Event{Data: userResponse,
				Topic: TopicLoggedInOut,
			})
			updateBookingsResponse.Error = UnauthorizedErrorStr
			return
		}
		if err != nil {
			if errors.Is(err, token.ErrExpiredToken) || errors.Is(err, token.ErrInvalidToken) {
				s.eventBroker.Fire(Event{Data: UserResponse{}, Topic: TopicLoggedInOut})
			}
			updateBookingsResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&updateBookingsResponse)
		if err != nil {
			updateBookingsResponse.Error = err.Error()
			return
		}
	}()
}

func (s *service) DeleteBooking(bookingID primitive.ObjectID) {
	go func() {
		var deleteBookingResponse DeleteBookingResponse
		var err error
		deleteBookingResponse.ID = bookingID
		if err != nil {
			alog.Logger().Println(err)
		}
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				deleteBookingResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  deleteBookingResponse,
				Topic: TopicDeleteBooking,
			})
		}()
		if bookingID.Hex() == primitive.NilObjectID.Hex() {
			deleteBookingResponse.Error = "booking id cannot be empty"
			return
		}
		userEvent, ok := s.eventBroker.cachedEvents.Get(TopicLoggedInOut)
		if !ok {
			deleteBookingResponse.Error = "user not logged in"
			return
		}
		userResponse, ok := userEvent.Data.(UserResponse)
		if !ok {
			deleteBookingResponse.Error = "critical error, need to contact admin"
			return
		}
		if userResponse.AccessToken == "" {
			deleteBookingResponse.Error = "user not logged in"
			return
		}
		delReq := DeleteBookingRequest{ID: bookingID}
		jsonValues, err := json.Marshal(delReq)
		if err != nil {
			deleteBookingResponse.Error = err.Error()
			return
		}
		req, err := http.NewRequest("DELETE", s.config.ApiURL+"/bookings", bytes.NewBuffer(jsonValues))
		if err != nil {
			deleteBookingResponse.Error = err.Error()
			return
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userResponse.AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if rsp.StatusCode == http.StatusUnauthorized {
			userResponse = UserResponse{AccessToken: "", User: response.User{}, Error: ""}
			s.eventBroker.Fire(Event{Data: userResponse,
				Topic: TopicLoggedInOut,
			})
			deleteBookingResponse.Error = UnauthorizedErrorStr
			return
		}
		if err != nil {
			if errors.Is(err, token.ErrExpiredToken) || errors.Is(err, token.ErrInvalidToken) {
				s.eventBroker.Fire(Event{Data: UserResponse{}, Topic: TopicLoggedInOut})
			}
			deleteBookingResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&deleteBookingResponse)
		if err != nil {
			deleteBookingResponse.Error = err.Error()
			return
		}
	}()
}
