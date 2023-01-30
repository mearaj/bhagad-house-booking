package request

import "github.com/mearaj/bhagad-house-booking/common/model"

type CreateUser struct {
	Name     string            `json:"name" bson:"name"`
	Email    string            `json:"email" bson:"email" binding:"required,email"`
	Password string            `json:"password" bson:"password" binding:"required,min=6"`
	Roles    []model.UserRoles `json:"roles" bson:"roles"`
}

type ListUsers struct {
	Limit  int32 `json:"limit" bson:"limit"`
	Offset int32 `json:"offset" bson:"offset"`
}

type UpdateUser struct {
	ID            int64  `json:"id" bson:"id"`
	Name          string `json:"name" bson:"name"`
	Email         string `json:"email" bson:"email"`
	EmailVerified bool   `json:"email_verified,omitempty" bson:"email_verified,omitempty"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
