package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transaction struct {
	ID               primitive.ObjectID `json:"_id,omitempty"        bson:"_id,omitempty"`
	BookingNumber    int                `json:"booking_number"        bson:"booking_number"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
	Amount           float64            `json:"amount" bson:"amount"`
	Details          string             `json:"details"    bson:"details"`
	PaymentMode      PaymentMode        `json:"payment_mode" bson:"payment_mode"`
	Number           int                `json:"number" bson:"number"`
	ConfirmEmailSent bool               `json:"confirm_email_sent" bson:"confirm_email_sent"`
	ConfirmSMSSent   bool               `json:"confirm_sms_sent" bson:"confirm_sms_sent"`
}

type PaymentMode int

const (
	PaymentModeCash = PaymentMode(iota)
	PaymentModeCheque
	PaymentModeCreditCard
	PaymentModeUnknown = -1
)

type PaymentModeString string

const (
	PaymentModeCashStr       PaymentModeString = "Cash"
	PaymentModeChequeStr     PaymentModeString = "Cheque"
	PaymentModeCreditCardStr PaymentModeString = "Credit Card"
	PaymentModeUnknownStr    PaymentModeString = "Unknown"
)

func (p PaymentMode) ModeString() PaymentModeString {
	switch p {
	case PaymentModeCash:
		return PaymentModeCashStr
	case PaymentModeCheque:
		return PaymentModeChequeStr
	case PaymentModeCreditCard:
		return PaymentModeCreditCardStr
	}
	return PaymentModeUnknownStr
}
func (p PaymentMode) String() string {
	return string(p.ModeString())
}

func (p PaymentModeString) ModeInt() PaymentMode {
	switch p {
	case PaymentModeCashStr:
		return PaymentModeCash
	case PaymentModeChequeStr:
		return PaymentModeCheque
	case PaymentModeCreditCardStr:
		return PaymentModeCreditCard
	}
	return PaymentModeUnknown
}
