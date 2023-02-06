package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	ID               primitive.ObjectID `json:"_id,omitempty"        bson:"_id,omitempty"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
	StartDate        time.Time          `json:"start_date" bson:"start_date"`
	EndDate          time.Time          `json:"end_date"   bson:"end_date"`
	Details          string             `json:"details"    bson:"details"`
	CustomerName     string             `json:"customer_name" bson:"customer_name"`
	RatePerDay       float64            `json:"rate_per_day" bson:"rate_per_day"`
	Number           int                `json:"number" bson:"number"`
	CustomerPhone    string             `json:"customer_phone" bson:"customer_phone"`
	CustomerEmail    string             `json:"customer_email" bson:"customer_email"`
	ConfirmEmailSent bool               `json:"confirm_email_sent" bson:"confirm_email_sent"`
	ConfirmSMSSent   bool               `json:"confirm_sms_sent" bson:"confirm_sms_sent"`
}
