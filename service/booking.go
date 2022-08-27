package service

import (
	"gorm.io/gorm"
	"time"
)

type Booking struct {
	gorm.Model
	StartDate time.Time
	EndDate   time.Time
	Rate      float64
	RateUnit  RateUnit
}
type RateUnit int

const (
	RatePerDay = RateUnit(iota)
	RatePerWeek
	RatePerMonth
	RatePerHour
)
