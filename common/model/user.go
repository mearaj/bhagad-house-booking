package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserRoles string

const (
	UserRolesUser  UserRoles = "User"
	UserRolesAdmin UserRoles = "Admin"
)

type User struct {
	EmailVerified     bool               `json:"email_verified,omitempty" bson:"email_verified,omitempty"`
	PasswordChangedAt time.Time          `json:"password_changed_at,omitempty" bson:"password_changed_at,omitempty"`
	CreatedAt         time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Password          string             `json:"password,omitempty" bson:"password,omitempty"`
	Name              string             `json:"name,omitempty" bson:"name,omitempty"`
	Email             string             `json:"email" bson:"email"`
	ID                primitive.ObjectID `json:"_id,omitempty"        bson:"_id,omitempty"`
	Roles             []UserRoles        `json:"roles" bson:"roles"`
}
