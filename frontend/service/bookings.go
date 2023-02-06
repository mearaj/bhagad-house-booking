package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mearaj/bhagad-house-booking/common/alog"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"net/http"
	"net/url"
)

const UnauthorizedErrorStr = "unauthorized"

func (s *service) GetBookings(bookingParams BookingsRequest) {
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if err != nil {
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		defer func() {
			_ = rsp.Body.Close()
		}()
		if err != nil {
			isAuthErr := s.FireAuthError(rsp, nil, err)
			if isAuthErr {
				searchBookingsResponse.Error = err.Error()
			}
			return
		}
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		defer func() {
			_ = rsp.Body.Close()
		}()
		if err != nil {
			isAuthErr := s.FireAuthError(rsp, nil, err)
			if isAuthErr {
				createBookingResponse.Error = UnauthorizedErrorStr
			}
			return
		}
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
		updateBookingsResponse.Booking.Number = bookingParams.Number
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		defer func() {
			_ = rsp.Body.Close()
		}()
		if err != nil {
			isAuthErr := s.FireAuthError(rsp, nil, err)
			if isAuthErr {
				updateBookingsResponse.Error = UnauthorizedErrorStr
			}
			return
		}
		err = json.NewDecoder(rsp.Body).Decode(&updateBookingsResponse)
		if err != nil {
			updateBookingsResponse.Error = err.Error()
			return
		}
	}()
}

func (s *service) DeleteBooking(number int) {
	go func() {
		var deleteBookingResponse DeleteBookingResponse
		var err error
		deleteBookingResponse.Number = number
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
		if number == 0 {
			deleteBookingResponse.Error = "booking id cannot zero"
			return
		}
		delReq := DeleteBookingRequest{Number: number}
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		defer func() {
			_ = rsp.Body.Close()
		}()
		if err != nil {
			isAuthErr := s.FireAuthError(rsp, nil, err)
			if isAuthErr {
				deleteBookingResponse.Error = UnauthorizedErrorStr
			}
			return
		}
		err = json.NewDecoder(rsp.Body).Decode(&deleteBookingResponse)
		if err != nil {
			deleteBookingResponse.Error = err.Error()
			return
		}
	}()
}
