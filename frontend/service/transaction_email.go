package service

import (
	"encoding/json"
	"fmt"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"net/http"
)

func (s *service) SendNewTransactionEmail(transactionNumber int, id interface{}) {
	go func() {
		var transactionEmailResponse NewTransactionEmailResponse
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				transactionEmailResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  transactionEmailResponse,
				Topic: TopicSendNewTransactionEmail,
				ID:    id,
			})
		}()
		req, err := http.NewRequest("POST", s.config.ApiURL+"/transactions/"+fmt.Sprintf("%d", transactionNumber)+"/sendNewTransactionEmail", nil)
		if err != nil {
			transactionEmailResponse.Error = err.Error()
			return
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if err != nil {
			isAuthErr := s.FireAuthError(rsp, id, err)
			if isAuthErr {
				transactionEmailResponse.Error = UnauthorizedErrorStr
			}
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&transactionEmailResponse)
		if err != nil {
			transactionEmailResponse.Error = err.Error()
			return
		}
	}()
}
