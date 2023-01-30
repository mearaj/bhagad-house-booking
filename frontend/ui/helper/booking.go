package helper

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"strconv"
	"strings"
	"time"
)

func IsNilObjectID(id primitive.ObjectID) bool {
	return id.Hex() == primitive.NilObjectID.Hex()
}

func BookingTotalPriceFromRateStr(bookingRatePerDay string, startDate, endDate time.Time) float64 {
	ratePerDayStr := strings.TrimSpace(bookingRatePerDay)
	if ratePerDayStr == "" {
		ratePerDayStr = "0"
	}
	ratePerDay, _ := strconv.ParseFloat(ratePerDayStr, 64)
	return BookingTotalPrice(ratePerDay, startDate, endDate)
}

func BookingTotalPrice(bookingRatePerDay float64, startDate, endDate time.Time) float64 {
	// if startDate and endDate are same, then it's 1 Day i.e. 24hrs, hence we need to add it EndDate
	dur := endDate.Add(time.Hour * 24).Sub(startDate)
	numberOfDays := float64(int(dur.Hours() / 24))
	totalPrice := math.Round(numberOfDays * bookingRatePerDay)
	return totalPrice
}
