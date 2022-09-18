package service

import (
	"time"
)

type Booking struct {
	ID        uint
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
