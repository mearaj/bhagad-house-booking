package utils

import (
	"github.com/mearaj/bhagad-house-booking/common/request"
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

func GetDefaultBookingRequest() request.ListBookings {
	startDate := GetFirstDayOfMonth(time.Now().Local())
	endDate := GetLastDayOfMonth(time.Now().Local().AddDate(0, 5, 0))
	return request.ListBookings{
		StartDate: startDate,
		EndDate:   endDate,
	}
}
func BookingTotalNumberOfDays(startDate, endDate time.Time) int {
	// if startDate and endDate are same, then it's 1 Day i.e. 24hrs, hence we need to add it EndDate
	startDate, _ = time.Parse("2006-01-02", startDate.Format("2006-01-02"))
	endDate, _ = time.Parse("2006-01-02", endDate.Format("2006-01-02"))
	dur := endDate.Add(time.Hour * 24).Sub(startDate)
	numberOfDays := int(dur.Hours() / 24)
	return numberOfDays
}
