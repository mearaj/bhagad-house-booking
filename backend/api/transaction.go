package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/common/request"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

const ErrorInvalidBookingID = "invalid booking id"

func (s *Server) addUpdateTransaction(ctx *gin.Context) {
	var rq request.AddUpdateTransaction
	var rsp response.AddUpdateTransaction
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	if rq.BookingID.Hex() == primitive.NilObjectID.Hex() {
		rsp.Error = ErrorInvalidBookingID
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	var booking model.Booking
	err := bookingsCollection.FindOne(context.TODO(), bson.M{"_id": rq.BookingID}).Decode(&booking)
	if err != nil {
		rsp.Error = err.Error()
		if errors.Is(err, mongo.ErrNoDocuments) {
			rsp.Error = ErrorInvalidBookingID
		}
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	isNew := rq.ID.Hex() == primitive.NilObjectID.Hex()
	if isNew {
		rq.CreatedAt = time.Now()
		rq.UpdatedAt = time.Now()
		result, err := transactionsCollection.InsertOne(context.TODO(), rq)
		if err != nil {
			rsp.Error = err.Error()
			ctx.JSON(http.StatusInternalServerError, rsp)
			return
		}
		rq.ID = result.InsertedID.(primitive.ObjectID)
	}
	err = transactionsCollection.FindOne(context.TODO(), bson.M{"_id": rq.ID, "booking_id": rq.BookingID}).Decode(&rsp.Transaction)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	if !isNew {
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "updated_at", Value: time.Now()},
			{Key: "amount", Value: rq.Amount},
			{Key: "details", Value: rq.Details},
		}}}
		filters := bson.D{{Key: "_id", Value: rq.ID}, {Key: "booking_id", Value: rq.BookingID}}
		result, err := transactionsCollection.UpdateOne(context.TODO(), filters, update)
		if err != nil {
			rsp.Error = err.Error()
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		if result.ModifiedCount != 1 {
			rsp.Error = "unexpected error"
			ctx.JSON(http.StatusInternalServerError, rsp)
			return
		}
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (s *Server) getTransactions(ctx *gin.Context) {
	var rsp response.GetTransactions
	var rq request.GetTransactions
	bookingID := ctx.Param("booking_id")
	if bookingID == "" {
		rsp.Error = ErrorInvalidBookingID
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	if bookingID == primitive.NilObjectID.Hex() || !primitive.IsValidObjectID(bookingID) {
		rsp.Error = ErrorInvalidBookingID
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	var err error
	rq.BookingID, err = primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		rsp.Error = ErrorInvalidBookingID
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	cursor, err := transactionsCollection.Find(context.TODO(), bson.M{"booking_id": rq.BookingID})
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	if err = cursor.All(context.TODO(), &rsp.Transactions); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}
func (s *Server) deleteTransaction(ctx *gin.Context) {
	var rq request.DeleteTransaction
	var rsp response.DeleteTransaction
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	filter := bson.D{{Key: "_id", Value: rq.ID}, {Key: "booking_id", Value: rq.BookingID}}
	result, err := transactionsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if result.DeletedCount == 1 {
		rsp.ID = rq.ID
		rsp.BookingID = rq.BookingID
		ctx.JSON(http.StatusOK, rsp)
		return
	}
	rsp.Error = "no matched transaction found..."
	ctx.JSON(http.StatusBadRequest, rsp)
}
