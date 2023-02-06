package service

import (
	"encoding/json"
	"fmt"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"net/http"
)

func (s *service) SendNewBookingSMS(bookingNumber int, id interface{}) {
	go func() {
		var bookingSMSResponse NewBookingSMSResponse
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				bookingSMSResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  bookingSMSResponse,
				Topic: TopicSendNewBookingSMS,
				ID:    id,
			})
		}()
		req, err := http.NewRequest("POST", s.config.ApiURL+"/bookings/"+fmt.Sprintf("%d", bookingNumber)+"/sendNewBookingSMS", nil)
		if err != nil {
			bookingSMSResponse.Error = err.Error()
			return
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if err != nil {
			isAuthErr := s.FireAuthError(rsp, id, err)
			if isAuthErr {
				bookingSMSResponse.Error = UnauthorizedErrorStr
			}
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&bookingSMSResponse)
		if err != nil {
			bookingSMSResponse.Error = err.Error()
			return
		}
	}()
}
