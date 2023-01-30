package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"net/url"
	"strings"
	"time"
)

func validateGetBookingsQuery(ctx *gin.Context) (startTime, endTime time.Time, err error) {
	startTimeStr, err := url.QueryUnescape(ctx.Query("start_date"))
	if strings.TrimSpace(startTimeStr) == "" || err != nil {
		if err == nil {
			err = errors.New("invalid format, supported format yyyy-mm-dd")
		}
		return
	}
	endTimeStr, err := url.QueryUnescape(ctx.Query("end_date"))
	if strings.TrimSpace(endTimeStr) == "" || err != nil {
		if err == nil {
			err = errors.New("invalid format, supported format yyyy-mm-dd")
		}
		return
	}

	startTime, err = time.Parse("2006-01-02", startTimeStr)
	if err != nil {
		return
	}
	endTime, err = time.Parse("2006-01-02", endTimeStr)
	if err != nil {
		return
	}
	if endTime.Before(startTime) {
		err = errors.New("end time should be equal or after start time")
		return
	}
	return
}

func findConflictingBookings(startDate time.Time, endDate time.Time) (*[]model.Booking, error) {
	var bookings []model.Booking
	// Check if the booking conflicts with any current booking(s)
	cursor, err := bookingsCollection.Find(context.TODO(),
		bson.D{{Key: "$or", Value: bson.A{
			bson.D{{Key: "start_date", Value: bson.D{{Key: "$gte", Value: startDate}}},
				{Key: "start_date", Value: bson.D{{Key: "$lte", Value: endDate}}},
				{Key: "end_date", Value: bson.D{{Key: "$gte", Value: startDate}}},
				{Key: "end_date", Value: bson.D{{Key: "$gte", Value: endDate}}},
			},
			bson.D{{Key: "start_date", Value: bson.D{{Key: "$lte", Value: startDate}}},
				{Key: "start_date", Value: bson.D{{Key: "$lte", Value: endDate}}},
				{Key: "end_date", Value: bson.D{{Key: "$gte", Value: startDate}}},
				{Key: "end_date", Value: bson.D{{Key: "$lte", Value: endDate}}},
			},
			bson.D{{Key: "start_date", Value: bson.D{{Key: "$lte", Value: startDate}}},
				{Key: "start_date", Value: bson.D{{Key: "$lte", Value: endDate}}},
				{Key: "end_date", Value: bson.D{{Key: "$gte", Value: startDate}}},
				{Key: "end_date", Value: bson.D{{Key: "$gte", Value: endDate}}},
			},
			bson.D{{Key: "start_date", Value: bson.D{{Key: "$gte", Value: startDate}}},
				{Key: "start_date", Value: bson.D{{Key: "$lte", Value: endDate}}},
				{Key: "end_date", Value: bson.D{{Key: "$gte", Value: startDate}}},
				{Key: "end_date", Value: bson.D{{Key: "$lte", Value: endDate}}},
			},
			bson.D{{Key: "start_date", Value: bson.D{{Key: "$eq", Value: startDate}}},
				{Key: "end_date", Value: bson.D{{Key: "$eq", Value: endDate}}},
			},
		}}},
	)
	if err != nil {
		return &bookings, err
	}
	if err = cursor.All(context.TODO(), &bookings); err != nil {
		return &bookings, err
	}
	return &bookings, nil
}
