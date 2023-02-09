package service

import (
	"encoding/json"
	"fmt"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"net/http"
)

func (s *service) SendNewTransactionSMS(transactionNumber int, id interface{}) {
	go func() {
		var transactionSMSResponse NewTransactionSMSResponse
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				transactionSMSResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  transactionSMSResponse,
				Topic: TopicSendNewTransactionSMS,
				ID:    id,
			})
		}()
		req, err := http.NewRequest("POST", s.config.ApiURL+"/transactions/"+fmt.Sprintf("%d", transactionNumber)+"/sendNewTransactionSMS", nil)
		if err != nil {
			transactionSMSResponse.Error = err.Error()
			return
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if err != nil {
			isAuthErr := s.FireAuthError(rsp, id, err)
			if isAuthErr {
				transactionSMSResponse.Error = UnauthorizedErrorStr
			}
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&transactionSMSResponse)
		if err != nil {
			transactionSMSResponse.Error = err.Error()
			return
		}
	}()
}
