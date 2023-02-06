package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/common/request"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/url"
	"time"
)

func (s *Server) createBooking(ctx *gin.Context) {
	var rq request.CreateBooking
	var rsp response.CreateBooking
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	startDate, err := utils.GetFormatted20060102(rq.StartDate)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	endDate, err := utils.GetFormatted20060102(rq.EndDate)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	if endDate.Before(startDate) {
		rsp.Error = "end time should be equal or after start time"
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	bookings, err := findConflictingBookings(startDate, endDate)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if len(*bookings) > 0 {
		rsp.Error = "booking conflicts"
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	var findBooking model.Booking
	err = bookingsCollection.FindOne(context.TODO(), bson.D{}, &options.FindOneOptions{
		Sort: bson.D{{Key: "number", Value: -1}},
	}).Decode(&findBooking)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rq.Number = findBooking.Number + 1
	rq.StartDate = startDate
	rq.EndDate = endDate
	rq.CreatedAt = time.Now()
	result, err := bookingsCollection.InsertOne(context.TODO(), rq)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	id := result.InsertedID.(primitive.ObjectID)
	err = bookingsCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&rsp.Booking)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (s *Server) getBookings(ctx *gin.Context) {
	var rsp response.Bookings

	startTime, endTime, err := validateGetBookingsQuery(ctx)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	filter := bson.D{
		{Key: "start_date", Value: bson.D{{Key: "$gte", Value: startTime}}},
		{Key: "end_date", Value: bson.D{{Key: "$lte", Value: endTime}}},
	}
	cursor, err := bookingsCollection.Find(context.TODO(), filter)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if err = cursor.All(context.TODO(), &rsp.Bookings); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	// If the caller is unauthorized user, then details is empty
	if _, _, err = checkUserAuthorized(ctx, s.tokenMaker); err != nil {
		var anonymousBookings = make([]model.Booking, 0, 1)
		for _, booking := range rsp.Bookings {
			booking.Details = ""
			booking.RatePerDay = 0
			booking.CustomerName = ""
			anonymousBookings = append(anonymousBookings, booking)
			fmt.Printf("%+v\n", booking)
		}
		rsp.Bookings = anonymousBookings
		ctx.JSON(http.StatusOK, rsp)
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (s *Server) updateBooking(ctx *gin.Context) {
	var rq request.UpdateBooking
	var rsp response.UpdateBooking
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	rsp.Booking.Number = rq.Number
	startDateStr := rq.StartDate.Format("2006-01-02")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	endDateStr := rq.EndDate.Format("2006-01-02")
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	bookings, err := findConflictingBookings(startDate, endDate)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if len(*bookings) > 0 {
		// There should only be one booking and it should be booking request to be updated
		for _, booking := range *bookings {
			if booking.Number != rq.Number {
				rsp.Error = "booking conflicts"
				ctx.JSON(http.StatusBadRequest, rsp)
				return
			}
		}
	}

	rq.StartDate = startDate
	rq.EndDate = endDate
	filter := bson.D{{Key: "number", Value: rq.Number}}
	update := bson.D{{Key: "$set", Value: rq}}
	result, err := bookingsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if result.MatchedCount == 1 {
		err = bookingsCollection.FindOne(context.TODO(), bson.M{"number": rq.Number}).Decode(&rsp.Booking)
		if err != nil {
			rsp.Error = err.Error()
			ctx.JSON(http.StatusInternalServerError, rsp)
			return
		}
		ctx.JSON(http.StatusOK, rsp)
		return
	}
	rsp.Error = "no matched booking found..."
	ctx.JSON(http.StatusBadRequest, rsp)
}

func (s *Server) deleteBooking(ctx *gin.Context) {
	var rq request.DeleteBooking
	var rsp response.DeleteBooking
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	filter := bson.D{{Key: "number", Value: rq.Number}}
	result, err := bookingsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if result.DeletedCount == 1 {
		rsp.Number = rq.Number
		ctx.JSON(http.StatusOK, rsp)
		return
	}
	rsp.Error = "no matched booking found..."
	ctx.JSON(http.StatusBadRequest, rsp)
}

func (s *Server) searchBookings(ctx *gin.Context) {
	query, err := url.QueryUnescape(ctx.Query("query"))
	var rsp response.SearchBookings
	if err != nil {
		rsp.Error = "invalid search query"
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "customer_name", Value: bson.D{{Key: "$regex", Value: primitive.Regex{
				Pattern: query,
				Options: "i",
			}}}}},
			bson.D{{Key: "details", Value: bson.D{{Key: "$regex", Value: primitive.Regex{
				Pattern: query,
				Options: "i",
			}}}}},
		}},
	}
	cursor, err := bookingsCollection.Find(context.TODO(), filter)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if err = cursor.All(context.TODO(), &rsp.Bookings); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	// If the caller is unauthorized user, then details is empty
	if _, _, err = checkUserAuthorized(ctx, s.tokenMaker); err != nil {
		var anonymousBookings = make([]model.Booking, 0, 1)
		for _, booking := range rsp.Bookings {
			booking.Details = ""
			booking.RatePerDay = 0
			booking.CustomerName = ""
			anonymousBookings = append(anonymousBookings, booking)
			fmt.Printf("%+v\n", booking)
		}
		rsp.Bookings = anonymousBookings
		ctx.JSON(http.StatusOK, rsp)
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}
