package api

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/sendgrid/sendgrid-go"
	gridmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/mail"
	"strconv"
)

//go:embed email.html

var EmailHTML string

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

	from := gridmail.NewEmail("Bhagad House", s.config.AdminEmail)
	subject := "Bhagad House Booking Confirmation"
	to := gridmail.NewEmail(rsp.Booking.CustomerName, rsp.Booking.CustomerEmail)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := fmt.Sprintf(EmailHTML)
	message := gridmail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.config.SendGridAPIKey)
	gridResp, err := client.Send(message)
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
