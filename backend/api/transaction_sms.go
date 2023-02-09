package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

func (s *Server) sendNewTransactionSMS(ctx *gin.Context) {
	var rsp response.NewTransactionSMS
	transactionNumber := ctx.Param("number")
	if transactionNumber == "" {
		rsp.Error = ErrorInvalidTransactionNumber
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	var err error
	number, err := strconv.ParseInt(transactionNumber, 10, 64)
	if err != nil {
		rsp.Error = ErrorInvalidTransactionNumber
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	rsp.Transaction.Number = int(number)

	err = transactionsCollection.FindOne(context.TODO(), bson.D{{Key: "number", Value: number}}).Decode(&rsp.Transaction)
	if err != nil {
		rsp.Error = err.Error()
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	var foundBooking model.Booking
	err = bookingsCollection.FindOne(context.TODO(), bson.D{{Key: "number", Value: rsp.Transaction.BookingNumber}}).Decode(&foundBooking)
	if err != nil {
		rsp.Error = err.Error()
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	phoneNumber := foundBooking.CustomerPhone
	isNumberValid := utils.ValidateIndianPhoneNumber(phoneNumber)
	if !isNumberValid {
		rsp.Error = ErrorInvalidPhoneNumber
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	var transactions []model.Transaction
	cursor, err := transactionsCollection.Find(context.TODO(), bson.D{{Key: "booking_number", Value: foundBooking.Number}})
	if err != nil {
		rsp.Error = err.Error()
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if err = cursor.All(context.TODO(), &transactions); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	bookingPeriodInt := utils.BookingTotalNumberOfDays(foundBooking.StartDate, foundBooking.EndDate)
	bookingTotalPrice := foundBooking.RatePerDay * float64(bookingPeriodInt)

	var totalAmountReceived float64
	for _, trans := range transactions {
		totalAmountReceived += trans.Amount
	}

	previousBal := bookingTotalPrice - (totalAmountReceived - rsp.Transaction.Amount)
	balanceLeft := previousBal - rsp.Transaction.Amount

	switch len(phoneNumber) {
	case 10: // if there is no +91
		phoneNumber = "+91" + phoneNumber
	case 11: // if the number begins with 0
		phoneNumber = "+91" + phoneNumber[1:]
	case 12: // if the number begins without +
		phoneNumber = "+" + phoneNumber
	}
	bookingNumberStr := fmt.Sprintf("%d.", foundBooking.Number)
	bookingTotalPriceStr := fmt.Sprintf("INR %.2f.", bookingTotalPrice)
	textContent := fmt.Sprintf(
		"%s\n\n%s\n\n%s : %s\n%s : %s\n%s : %s\n%s : %s\n%s : %.2f\n%s : %.2f\n\n%s",
		"Payment received for your booking at https://bhagadhouse.com.",
		"Payment Details",
		"Booking No.", bookingNumberStr,
		"Receipt No.", bookingNumberStr,
		"Payment Mode", rsp.Transaction.PaymentMode,
		"Total Price", bookingTotalPriceStr,
		"Amount Received", rsp.Transaction.Amount,
		"Balance Left", balanceLeft,
		"Than you for your business",
	)
	// Find your Account SID and Auth Token at twilio.com/console
	// and set the environment variables. See http://twil.io/secure
	client := twilio.NewRestClient()

	params := &api.CreateMessageParams{}
	params.SetBody(textContent)
	params.SetFrom(s.config.TwilioPhoneNumber)
	params.SetTo(phoneNumber)

	_, err = client.Api.CreateMessage(params)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	filter := bson.D{{Key: "number", Value: rsp.Transaction.Number}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"confirm_sms_sent": true}}}
	result, err := transactionsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if result.MatchedCount == 1 {
		rsp.Transaction.ConfirmSMSSent = true
	} else {
		rsp.Error = "unexpected error"
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
