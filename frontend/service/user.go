package service

import (
	"encoding/json"
	"fmt"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"net/http"
	"strings"
)

type UserResponse sqlc.LoginUserResponse

func (s *UserResponse) IsLoggedIn() bool {
	return s.AccessToken != ""
}

func (s *UserResponse) IsAdmin() (isAdmin bool) {
	if s.AccessToken == "" {
		return isAdmin
	}
	for _, role := range s.User.Roles {
		if role == sqlc.UserRolesAdmin {
			isAdmin = true
			break
		}
	}
	return isAdmin
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *service) LogInUser(email string, password string) {
	go func() {
		var resp UserResponse
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from err, ", r)
				resp.Error = fmt.Sprintf("%v", r)
			}
			s.eventBroker.Fire(Event{
				Data:  resp,
				Topic: TopicUserLoggedInOut,
			})
		}()
		loginReq := LoginRequest{Email: email, Password: password}
		loginReqMar, err := json.Marshal(loginReq)
		if err != nil {
			resp.Error = err.Error()
			return
		}
		req, err := http.NewRequest("POST", s.config.ApiURL+"/users/login", strings.NewReader(string(loginReqMar)))
		if err != nil {
			resp.Error = err.Error()
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Accept", "application/json")
		cl := http.Client{}
		httpResp, err := cl.Do(req)
		if err != nil {
			resp.Error = err.Error()
			return
		}
		defer func() {
			_ = httpResp.Body.Close()
		}()
		err = json.NewDecoder(httpResp.Body).Decode(&resp)
		if err != nil {
			resp.Error = err.Error()
		}
	}()
}

func (s *service) LogOutUser() {
	go func() {
		defer func() {
		}()
		s.eventBroker.Fire(Event{
			Data:  UserResponse{},
			Topic: TopicUserLoggedInOut,
		})
	}()
}
