package utils

import (
	"fmt"
	"time"
)

func GetFormattedDate(t time.Time) string {
	var bookingDate string
	switch t.Day() {
	case 1, 21, 31:
		bookingDate = fmt.Sprintf("%dst", t.Day())
	case 2, 22:
		bookingDate = fmt.Sprintf("%dnd", t.Day())
	case 3, 23:
		bookingDate = fmt.Sprintf("%drd", t.Day())
	default:
		bookingDate = fmt.Sprintf("%dth", t.Day())
	}
	month := t.Month().String()
	year := t.Year()
	day := t.Weekday().String()
	bookingDate = fmt.Sprintf("%s %s, %s, %d", bookingDate, month, day, year)
	return bookingDate
}
