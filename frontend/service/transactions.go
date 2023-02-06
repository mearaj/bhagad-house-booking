package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mearaj/bhagad-house-booking/common/alog"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (s *service) GetTransactions(transactionsParams TransactionsRequest, id interface{}) {
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
		req, err := http.NewRequest("GET", s.config.ApiURL+"/bookings/"+fmt.Sprintf("%d", transactionsParams.BookingNumber)+"/transactions", nil)
		if err != nil {
			transactionsResponse.Error = err.Error()
			return
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if err != nil {
			isAuthErr := s.FireAuthError(rsp, id, err)
			if isAuthErr {
				transactionsResponse.Error = UnauthorizedErrorStr
			}
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
		if request.BookingNumber == 0 {
			addUpdateResponse.Error = "booking id cannot be empty"
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if err != nil {
			isAuthErr := s.FireAuthError(rsp, eventId, err)
			if isAuthErr {
				addUpdateResponse.Error = UnauthorizedErrorStr
			}
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
		var delTransResponse DeleteTransactionResponse
		var err error
		if err != nil {
			alog.Logger().Println(err)
		}
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				delTransResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  delTransResponse,
				Topic: TopicDeleteTransaction,
				ID:    eventID,
			})
		}()
		if req.BookingNumber == 0 {
			delTransResponse.Error = "booking id cannot be empty"
			return
		}
		if req.ID.Hex() == primitive.NilObjectID.Hex() {
			delTransResponse.Error = "transaction id cannot be empty"
			return
		}

		jsonValues, err := json.Marshal(req)
		if err != nil {
			delTransResponse.Error = err.Error()
			return
		}
		httpReq, err := http.NewRequest("DELETE", s.config.ApiURL+"/transactions", bytes.NewBuffer(jsonValues))
		if err != nil {
			delTransResponse.Error = err.Error()
			return
		}
		httpReq.Header.Add("Accept", "application/json")
		httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(httpReq)
		if err != nil {
			isAuthErr := s.FireAuthError(rsp, eventID, err)
			if isAuthErr {
				delTransResponse.Error = UnauthorizedErrorStr
			}
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&delTransResponse)
		if err != nil {
			delTransResponse.Error = err.Error()
			return
		}
	}()
}
