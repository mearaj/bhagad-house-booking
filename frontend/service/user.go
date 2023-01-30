package service

import (
	"encoding/json"
	"fmt"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *service) LogInUser(email string, password string) {
	go func() {
		var rsp UserResponse
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				rsp.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  rsp,
				Topic: TopicLoggedInOut,
			})
		}()
		loginReq := LoginRequest{Email: email, Password: password}
		loginReqMar, err := json.Marshal(loginReq)
		if err != nil {
			rsp.Error = err.Error()
			return
		}
		req, err := http.NewRequest("POST", s.config.ApiURL+"/users/login", strings.NewReader(string(loginReqMar)))
		if err != nil {
			rsp.Error = err.Error()
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Accept", "application/json")
		cl := http.Client{}
		httpResp, err := cl.Do(req)
		if err != nil {
			rsp.Error = err.Error()
			return
		}
		defer func() {
			_ = httpResp.Body.Close()
		}()
		err = json.NewDecoder(httpResp.Body).Decode(&rsp)
		if err != nil {
			rsp.Error = err.Error()
		}
	}()
}

func (s *service) LogOutUser() {
	go func() {
		defer func() {
		}()
		s.eventBroker.Fire(Event{
			Data:  response.LoginUser{},
			Topic: TopicLoggedInOut,
		})
	}()
}
