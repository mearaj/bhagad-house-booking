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
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"time"
)

const ErrorInvalidBookingID = "invalid booking id"
const ErrorInvalidTransactionNumber = "invalid transaction number"
const ErrorInvalidPhoneNumber = "invalid phone number"

func (s *Server) addUpdateTransaction(ctx *gin.Context) {
	var rq request.AddUpdateTransaction
	var rsp response.AddUpdateTransaction
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	if rq.BookingNumber == 0 {
		rsp.Error = ErrorInvalidBookingID
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	var booking model.Booking
	err := bookingsCollection.FindOne(context.TODO(), bson.M{"number": rq.BookingNumber}).Decode(&booking)
	if err != nil {
		rsp.Error = err.Error()
		if errors.Is(err, mongo.ErrNoDocuments) {
			rsp.Error = ErrorInvalidBookingID
		}
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	isNew := rq.ID.Hex() == primitive.NilObjectID.Hex()

	// foundTransaction's purpose is to find the maximum value of transaction number for inserting
	var foundTransaction model.Transaction
	err = transactionsCollection.FindOne(context.TODO(), bson.D{}, &options.FindOneOptions{
		Sort: bson.D{{Key: "number", Value: -1}},
	}).Decode(&foundTransaction)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	if isNew {
		rq.CreatedAt = time.Now()
		rq.UpdatedAt = time.Now()
		rq.Number = foundTransaction.Number + 1
		result, err := transactionsCollection.InsertOne(context.TODO(), rq)
		if err != nil {
			rsp.Error = err.Error()
			ctx.JSON(http.StatusInternalServerError, rsp)
			return
		}
		rq.ID = result.InsertedID.(primitive.ObjectID)
	}
	err = transactionsCollection.FindOne(context.TODO(), bson.M{"_id": rq.ID, "booking_number": rq.BookingNumber}).Decode(&rsp.Transaction)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	if !isNew {
		value := bson.D{
			{Key: "updated_at", Value: time.Now()},
			{Key: "amount", Value: rq.Amount},
			{Key: "details", Value: rq.Details},
			{Key: "payment_mode", Value: rq.PaymentMode},
		}
		if rq.Number == 0 {
			rq.Number = foundTransaction.Number + 1
			value = append(value, bson.E{Key: "number", Value: rq.Number})
		}
		update := bson.D{{Key: "$set", Value: value}}
		filters := bson.D{{Key: "_id", Value: rq.ID}, {Key: "booking_number", Value: rq.BookingNumber}}
		result, err := transactionsCollection.UpdateOne(context.TODO(), filters, update)
		if err != nil {
			rsp.Error = err.Error()
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		if result.ModifiedCount != 1 {
			rsp.Error = UnexpectedError
			ctx.JSON(http.StatusInternalServerError, rsp)
			return
		}
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (s *Server) getTransactions(ctx *gin.Context) {
	var rsp response.GetTransactions
	var rq request.GetTransactions
	bookingID := ctx.Param("number")
	if bookingID == "" {
		rsp.Error = ErrorInvalidBookingID
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	var err error
	number, err := strconv.ParseInt(bookingID, 10, 64)
	rq.BookingNumber = int(number)
	if err != nil {
		rsp.Error = ErrorInvalidBookingID
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	cursor, err := transactionsCollection.Find(context.TODO(), bson.M{"booking_number": rq.BookingNumber})
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
	filter := bson.D{{Key: "_id", Value: rq.ID}, {Key: "booking_number", Value: rq.BookingNumber}}
	result, err := transactionsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if result.DeletedCount == 1 {
		rsp.ID = rq.ID
		rsp.BookingNumber = rq.BookingNumber
		ctx.JSON(http.StatusOK, rsp)
		return
	}
	rsp.Error = "no matched transaction found..."
	ctx.JSON(http.StatusBadRequest, rsp)
}
