package api

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/sendgrid/sendgrid-go"
	gridmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/mail"
	"strconv"
)

//go:embed booking_email.html
var BookingEmailHTML string

const UnexpectedError = "unexpected error"

func (s *Server) sendNewBookingEmail(ctx *gin.Context) {
	var rsp response.NewBookingEmail
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
	_, err = mail.ParseAddress(rsp.Booking.CustomerEmail)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	if s.config.SendGridAPIKey == "" {
		rsp.Error = UnexpectedError
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	from := gridmail.NewEmail("Bhagad House", s.config.AdminEmail)
	subject := "Bhagad House Booking Confirmation"
	to := gridmail.NewEmail(rsp.Booking.CustomerName, rsp.Booking.CustomerEmail)
	bookingNumberStr := fmt.Sprintf("%d", rsp.Booking.Number)
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
	htmlContent := fmt.Sprintf(
		BookingEmailHTML,
		bookingNumberStr,
		bookingStartStr,
		bookingEndStr,
		bookingPeriodStr,
		bookingRateStr,
		bookingTotalPriceStr,
	)
	message := gridmail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(s.config.SendGridAPIKey)
	gridResp, err := client.Send(message)
	if gridResp.StatusCode != http.StatusAccepted {
		rsp.Error = UnexpectedError
		ctx.JSON(gridResp.StatusCode, rsp)
		return
	}
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(gridResp.StatusCode, rsp)
		return
	}
	filter := bson.D{{Key: "number", Value: rsp.Booking.Number}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"confirm_email_sent": true}}}
	result, err := bookingsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	if result.MatchedCount == 1 {
		rsp.Booking.ConfirmEmailSent = true
	} else {
		rsp.Error = UnexpectedError
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}
