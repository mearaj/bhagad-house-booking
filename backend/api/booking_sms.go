package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

func (s *Server) sendNewBookingSMS(ctx *gin.Context) {
	var rsp response.NewBookingSMS
	bookingID := ctx.Param("number")
	if bookingID == "" {
		rsp.Error = ErrorInvalidBookingID
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	var err error
	number, err := strconv.ParseInt(bookingID, 10, 64)
	if err != nil {
		rsp.Error = ErrorInvalidBookingID
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	rsp.Booking.Number = int(number)

	err = bookingsCollection.FindOne(context.TODO(), bson.D{{Key: "number", Value: number}}).Decode(&rsp.Booking)
	if err != nil {
		rsp.Error = err.Error()
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	phoneNumber := rsp.Booking.CustomerPhone
	isNumberValid := utils.ValidateIndianPhoneNumber(phoneNumber)
	if !isNumberValid {
		rsp.Error = ErrorInvalidPhoneNumber
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	switch len(phoneNumber) {
	case 10: // if there is no +91
		phoneNumber = "+91" + phoneNumber
	case 11: // if the number begins with 0
		phoneNumber = "+91" + phoneNumber[1:]
	case 12: // if the number begins without +
		phoneNumber = "+" + phoneNumber
	}
	bookingNumberStr := fmt.Sprintf("%d.", rsp.Booking.Number)
	bookingStartStr := utils.GetFormattedDate(rsp.Booking.StartDate) + "."
	bookingEndStr := utils.GetFormattedDate(rsp.Booking.EndDate) + "."
	bookingPeriodInt := utils.BookingTotalNumberOfDays(rsp.Booking.StartDate, rsp.Booking.EndDate)
	bookingPeriodStr := fmt.Sprintf("%d day", bookingPeriodInt)
	if bookingPeriodInt > 1 {
		bookingPeriodStr = fmt.Sprintf("%d days", bookingPeriodInt)
	}
	bookingRateStr := fmt.Sprintf("%.2f", rsp.Booking.RatePerDay)
	bookingTotalPriceFloat := rsp.Booking.RatePerDay * float64(bookingPeriodInt)
	bookingTotalPriceStr := fmt.Sprintf("INR %.2f", bookingTotalPriceFloat)
	textContent := fmt.Sprintf(
		"%s\n%s\n%s : %s\n%s : %s\n%s : %s\n%s : %s\n%s : %s\n%s : %s\n%s : %s\n\n%s",
		"This is a confirmation sms for your booking at https://bhagadhouse.com.",
		"Booking Details",
		"Status", "Confirmed.",
		"Number", bookingNumberStr,
		"Starts From", bookingStartStr,
		"Ends On", bookingEndStr,
		"Period", bookingPeriodStr,
		"Rate", bookingRateStr,
		"Total Price", bookingTotalPriceStr,
		"Thank you for your business!",
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
	filter := bson.D{{Key: "number", Value: rsp.Booking.Number}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"confirm_sms_sent": true}}}
	result, err := bookingsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if result.MatchedCount == 1 {
		rsp.Booking.ConfirmSMSSent = true
	} else {
		rsp.Error = "unexpected error"
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
