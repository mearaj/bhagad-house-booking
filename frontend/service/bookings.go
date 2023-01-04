package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"net/http"
	"net/url"
)

const UnauthorizedErrorStr = "unauthorized"

func (s *service) Bookings(bookingParams BookingParams) {
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
				Topic: TopicBookingsFetched,
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
		userEvent, ok := s.eventBroker.cachedEvents.Value(TopicUserLoggedInOut)
		if ok {
			userResponse, ok := userEvent.Data.(UserResponse)
			if ok && userResponse.AccessToken != "" {
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userResponse.AccessToken))
			}
		}
		cl := http.Client{}
		resp, err := cl.Do(req)
		if err != nil {
			bookingsResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		err = json.NewDecoder(resp.Body).Decode(&bookingsResponse)
		if err != nil {
			bookingsResponse.Error = err.Error()
			return
		}
	}()
}
func (s *service) CreateBooking(bookingParams sqlc.CreateBookingParams) {
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
		userEvent, ok := s.eventBroker.cachedEvents.Value(TopicUserLoggedInOut)
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
		resp, err := cl.Do(req)
		if resp.StatusCode == http.StatusUnauthorized {
			s.eventBroker.Fire(Event{Data: UserResponse{AccessToken: "", User: sqlc.User{}, Error: ""},
				Topic: TopicUserLoggedInOut,
			})
			createBookingResponse.Error = UnauthorizedErrorStr
			return
		}
		if err != nil {
			createBookingResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		err = json.NewDecoder(resp.Body).Decode(&createBookingResponse)
		if err != nil {
			createBookingResponse.Error = err.Error()
			return
		}
	}()
}

func (s *service) UpdateBooking(bookingParams sqlc.UpdateBookingParams) {
	go func() {
		var updateBookingsResponse sqlc.UpdateBookingResponse
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
		userEvent, ok := s.eventBroker.cachedEvents.Value(TopicUserLoggedInOut)
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
		resp, err := cl.Do(req)
		if resp.StatusCode == http.StatusUnauthorized {
			s.eventBroker.Fire(Event{Data: UserResponse{AccessToken: "", User: sqlc.User{}, Error: ""},
				Topic: TopicUserLoggedInOut,
			})
			updateBookingsResponse.Error = UnauthorizedErrorStr
			return
		}
		if err != nil {
			updateBookingsResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		err = json.NewDecoder(resp.Body).Decode(&updateBookingsResponse)
		if err != nil {
			updateBookingsResponse.Error = err.Error()
			return
		}
	}()
}

type deleteBookingReq struct {
	ID int64 `json:"ID"`
}

func (s *service) DeleteBooking(bookingID int64) {
	go func() {
		var deleteBookingResponse DeleteBookingResponse
		deleteBookingResponse.ID = bookingID
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
		if bookingID == 0 {
			deleteBookingResponse.Error = "booking id cannot be 0"
			return
		}
		userEvent, ok := s.eventBroker.cachedEvents.Value(TopicUserLoggedInOut)
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
		delReq := deleteBookingReq{ID: bookingID}
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
		resp, err := cl.Do(req)
		if resp.StatusCode == http.StatusUnauthorized {
			s.eventBroker.Fire(Event{Data: UserResponse{AccessToken: "", User: sqlc.User{}, Error: ""},
				Topic: TopicUserLoggedInOut,
			})
			deleteBookingResponse.Error = UnauthorizedErrorStr
			return
		}
		if err != nil {
			deleteBookingResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		err = json.NewDecoder(resp.Body).Decode(&deleteBookingResponse)
		if err != nil {
			deleteBookingResponse.Error = err.Error()
			return
		}
	}()
}
