package response

import (
	"github.com/mearaj/bhagad-house-booking/common/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewUser struct {
	User  User   `json:"user,omitempty"`
	Error string `json:"error,omitempty"`
}

type User struct {
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Email string             `json:"email" bson:"email"`
	ID    primitive.ObjectID `json:"_id,omitempty"        bson:"_id,omitempty"`
	Roles []model.UserRoles  `json:"roles" bson:"roles"`
}

type LoginUser struct {
	AccessToken string `json:"access_token,omitempty"`
	User        User   `json:"user,omitempty"`
	Error       string `json:"error,omitempty"`
}

func (s *LoginUser) IsLoggedIn() bool {
	return s.AccessToken != ""
}

func (s *LoginUser) IsAdmin() (isAdmin bool) {
	if s.AccessToken == "" {
		return isAdmin
	}
	for _, role := range s.User.Roles {
		if role == model.UserRolesAdmin {
			isAdmin = true
			break
		}
	}
	return isAdmin
}

func (s *LoginUser) IsAuthorized() bool {
	return s.IsLoggedIn() && s.IsAdmin()
}

type Users struct {
	Users []User `json:"users,omitempty"`
	Error string `json:"error,omitempty"`
}
