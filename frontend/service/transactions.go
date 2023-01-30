package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mearaj/bhagad-house-booking/common/alog"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/common/token"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (s *service) Transactions(transactionsParams TransactionsRequest, id interface{}) {
	go func() {
		var transactionsResponse TransactionsResponse
		transactionsResponse.Transactions = make([]model.Transaction, 0)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				transactionsResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  transactionsResponse,
				Topic: TopicFetchTransactions,
				ID:    id,
			})
		}()
		req, err := http.NewRequest("GET", s.config.ApiURL+"/bookings/"+transactionsParams.BookingID.Hex()+"/transactions", nil)
		if err != nil {
			transactionsResponse.Error = err.Error()
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
		if rsp.StatusCode == http.StatusUnauthorized {
			userResponse := UserResponse{AccessToken: "", User: response.User{}, Error: ""}
			s.eventBroker.Fire(Event{Data: userResponse,
				Topic: TopicLoggedInOut,
			})
			userResponse.Error = UnauthorizedErrorStr
			*user.User() = userResponse
			return
		}
		if err != nil {
			if errors.Is(err, token.ErrExpiredToken) || errors.Is(err, token.ErrInvalidToken) {
				s.eventBroker.Fire(Event{Data: UserResponse{}, Topic: TopicLoggedInOut})
			}
			transactionsResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&transactionsResponse)
		if err != nil {
			transactionsResponse.Error = err.Error()
			return
		}
	}()
}
func (s *service) AddUpdateTransaction(request AddUpdateTransactionRequest, eventId interface{}) {
	go func() {
		var addUpdateResponse AddUpdateTransactionResponse
		var err error
		if err != nil {
			alog.Logger().Println(err)
		}
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				addUpdateResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  addUpdateResponse,
				Topic: TopicAddUpdateTransaction,
				ID:    eventId,
			})
		}()
		if request.BookingID.Hex() == primitive.NilObjectID.Hex() {
			addUpdateResponse.Error = "booking id cannot be empty"
			return
		}
		userEvent, ok := s.eventBroker.cachedEvents.Get(TopicLoggedInOut)
		if !ok {
			addUpdateResponse.Error = "user not logged in"
			return
		}
		userResponse, ok := userEvent.Data.(UserResponse)
		if !ok {
			addUpdateResponse.Error = "critical error, need to contact admin"
			return
		}
		if userResponse.AccessToken == "" {
			addUpdateResponse.Error = "user not logged in"
			return
		}
		jsonValues, err := json.Marshal(request)
		if err != nil {
			addUpdateResponse.Error = err.Error()
			return
		}
		req, err := http.NewRequest("POST", s.config.ApiURL+"/transactions", bytes.NewBuffer(jsonValues))
		if err != nil {
			addUpdateResponse.Error = err.Error()
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
			addUpdateResponse.Error = UnauthorizedErrorStr
			return
		}
		if err != nil {
			if errors.Is(err, token.ErrExpiredToken) || errors.Is(err, token.ErrInvalidToken) {
				s.eventBroker.Fire(Event{Data: UserResponse{}, Topic: TopicLoggedInOut})
			}
			addUpdateResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&addUpdateResponse)
		if err != nil {
			addUpdateResponse.Error = err.Error()
			return
		}
	}()
}
func (s *service) DeleteTransaction(req DeleteTransactionRequest, eventID interface{}) {
	go func() {
		var deleteTransactionResponse DeleteTransactionResponse
		var err error
		if err != nil {
			alog.Logger().Println(err)
		}
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				deleteTransactionResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  deleteTransactionResponse,
				Topic: TopicDeleteTransaction,
				ID:    eventID,
			})
		}()
		if req.BookingID.Hex() == primitive.NilObjectID.Hex() {
			deleteTransactionResponse.Error = "booking id cannot be empty"
			return
		}
		if req.ID.Hex() == primitive.NilObjectID.Hex() {
			deleteTransactionResponse.Error = "transaction id cannot be empty"
			return
		}
		userEvent, ok := s.eventBroker.cachedEvents.Get(TopicLoggedInOut)
		if !ok {
			deleteTransactionResponse.Error = "user not logged in"
			return
		}
		userResponse, ok := userEvent.Data.(UserResponse)
		if !ok {
			deleteTransactionResponse.Error = "critical error, need to contact admin"
			return
		}
		if userResponse.AccessToken == "" {
			deleteTransactionResponse.Error = "user not logged in"
			return
		}
		jsonValues, err := json.Marshal(req)
		if err != nil {
			deleteTransactionResponse.Error = err.Error()
			return
		}
		httpReq, err := http.NewRequest("DELETE", s.config.ApiURL+"/transactions", bytes.NewBuffer(jsonValues))
		if err != nil {
			deleteTransactionResponse.Error = err.Error()
			return
		}
		httpReq.Header.Add("Accept", "application/json")
		httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userResponse.AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(httpReq)
		if rsp.StatusCode == http.StatusUnauthorized {
			userResponse = UserResponse{AccessToken: "", User: response.User{}, Error: ""}
			s.eventBroker.Fire(Event{Data: userResponse,
				Topic: TopicLoggedInOut,
			})
			deleteTransactionResponse.Error = UnauthorizedErrorStr
			return
		}
		if err != nil {
			if errors.Is(err, token.ErrExpiredToken) || errors.Is(err, token.ErrInvalidToken) {
				s.eventBroker.Fire(Event{Data: UserResponse{}, Topic: TopicLoggedInOut})
			}
			deleteTransactionResponse.Error = err.Error()
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&deleteTransactionResponse)
		if err != nil {
			deleteTransactionResponse.Error = err.Error()
			return
		}
	}()
}
