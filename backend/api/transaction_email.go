package api

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/sendgrid/sendgrid-go"
	gridmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
)

//go:embed transaction_email.html
var TransactionEmailHTML string

func (s *Server) sendNewTransactionEmail(ctx *gin.Context) {
	var rsp response.NewTransactionEmail
	bookingID := ctx.Param("number")
	if bookingID == "" {
		rsp.Error = ErrorInvalidTransactionNumber
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	var err error
	number, err := strconv.ParseInt(bookingID, 10, 64)
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

	_, err = mail.ParseAddress(foundBooking.CustomerEmail)
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
	_ = balanceLeft

	from := gridmail.NewEmail("Bhagad House", s.config.AdminEmail)
	subject := "Bhagad House Payment Receipt"
	to := gridmail.NewEmail(foundBooking.CustomerName, foundBooking.CustomerEmail)
	bookingNumberStr := fmt.Sprintf(" %d ", foundBooking.Number)
	transactionNumberStr := fmt.Sprintf(" %d ", rsp.Transaction.Number)
	transactionDate := utils.GetFormattedDate(rsp.Transaction.CreatedAt)
	customerName := foundBooking.CustomerName
	paymentModeStr := fmt.Sprintf("%.2f <span style=\"font-weight:normal;\">through %s.</span>", rsp.Transaction.Amount,
		strings.ToLower(rsp.Transaction.PaymentMode.String()),
	)
	if strings.TrimSpace(customerName) == "" {
		customerName = foundBooking.CustomerEmail
	}

	htmlContent := fmt.Sprintf(
		TransactionEmailHTML,
		bookingNumberStr,
		transactionNumberStr,
		transactionDate,
		customerName,
		paymentModeStr,
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
	filter := bson.D{{Key: "number", Value: rsp.Transaction.Number}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"confirm_email_sent": true}}}
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
