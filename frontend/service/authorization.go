package service

import (
	"errors"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/common/token"
	"github.com/mearaj/bhagad-house-booking/frontend/user"
	"net/http"
	"strings"
)

func (s *service) FireAuthError(rsp *http.Response, id interface{}, err error) bool {
	isAuthError := errors.Is(err, token.ErrExpiredToken) ||
		errors.Is(err, token.ErrInvalidToken) ||
		strings.Contains(strings.ToLower(err.Error()), UnauthorizedErrorStr) ||
		rsp.StatusCode == http.StatusUnauthorized
	if isAuthError {
		userResponse := UserResponse{AccessToken: "", User: response.User{}, Error: ""}
		s.eventBroker.Fire(Event{Data: userResponse,
			Topic: TopicLoggedInOut,
			ID:    id,
		})
		userResponse.Error = UnauthorizedErrorStr
		user.SetUser(userResponse)
	}
	return isAuthError
}
