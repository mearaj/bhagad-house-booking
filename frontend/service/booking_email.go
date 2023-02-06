package service

import (
	"encoding/json"
	"fmt"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"net/http"
)

func (s *service) SendNewBookingEmail(bookingNumber int, id interface{}) {
	go func() {
		var bookingEmailResponse NewBookingEmailResponse
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				bookingEmailResponse.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  bookingEmailResponse,
				Topic: TopicSendNewBookingSMS,
				ID:    id,
			})
		}()
		req, err := http.NewRequest("POST", s.config.ApiURL+"/bookings/"+fmt.Sprintf("%d", bookingNumber)+"/sendNewBookingEmail", nil)
		if err != nil {
			bookingEmailResponse.Error = err.Error()
			return
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.User().AccessToken))
		cl := http.Client{}
		rsp, err := cl.Do(req)
		if err != nil {
			isAuthErr := s.FireAuthError(rsp, id, err)
			if isAuthErr {
				bookingEmailResponse.Error = UnauthorizedErrorStr
			}
			return
		}
		defer func() {
			_ = rsp.Body.Close()
		}()
		err = json.NewDecoder(rsp.Body).Decode(&bookingEmailResponse)
		if err != nil {
			bookingEmailResponse.Error = err.Error()
			return
		}
	}()
}
