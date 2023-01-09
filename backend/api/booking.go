package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"net/http"
	"net/url"
	"time"
)

type createBookingRequest struct {
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Details      string    `json:"details"`
	CustomerName string    `json:"customer_name"`
	TotalPrice   float64   `json:"total_price"`
}

func (s *Server) createBooking(ctx *gin.Context) {
	var req createBookingRequest
	var resp sqlc.CreateBookingResponse
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	startDateStr := req.StartDate.Format("2006-01-02")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	endDateStr := req.EndDate.Format("2006-01-02")
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	req.StartDate = startDate
	req.EndDate = endDate

	arg := sqlc.CreateBookingParams(req)
	booking, err := s.store.CreateBooking(ctx, arg)
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Booking = booking
	ctx.JSON(http.StatusOK, resp)
}

type getBookingRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getBooking(ctx *gin.Context) {
	//var req getBookingRequest
	//if err := ctx.ShouldBindUri(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	//	return
	//}
	//booking, err := s.store.GetBooking(ctx, req.ID)
	//if err != nil {
	//	if errors.Is(err, sql.ErrNoRows) {
	//		ctx.JSON(http.StatusNotFound, errorResponse(err))
	//	}
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//	return
	//}
	//ctx.JSON(http.StatusOK, booking)
}

func (s *Server) getBookings(ctx *gin.Context) {
	startTimeStr, err := url.QueryUnescape(ctx.Query("start_date"))
	var resp sqlc.BookingsResponse
	if startTimeStr == "" || err != nil {
		resp.Error = "invalid time format"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	endTimeStr, err := url.QueryUnescape(ctx.Query("end_date"))
	if endTimeStr == "" || err != nil {
		resp.Error = "invalid time format"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	startTime, err := time.Parse("2006-01-02", startTimeStr)
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	endTime, err := time.Parse("2006-01-02", endTimeStr)
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	bookings, err := s.store.ListBookings(ctx, sqlc.ListBookingsParams{
		StartDate: startTime,
		EndDate:   endTime,
	})
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	// If the caller is unauthorized user, then details is empty
	if _, _, err = checkUserAuthorized(ctx, s.tokenMaker); err != nil {
		var anonymousBookings []sqlc.Booking = make([]sqlc.Booking, 0, 1)
		for _, booking := range bookings {
			booking.Details = ""
			booking.TotalPrice = 0
			booking.CustomerName = ""
			anonymousBookings = append(anonymousBookings, booking)
			fmt.Printf("%+v\n", booking)
		}
		resp.Bookings = anonymousBookings
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.Bookings = bookings
	ctx.JSON(http.StatusOK, resp)
}

func (s *Server) updateBooking(ctx *gin.Context) {
	var req sqlc.UpdateBookingParams
	var resp sqlc.UpdateBookingResponse
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.Booking.ID = req.ID
	startDateStr := req.StartDate.Format("2006-01-02")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	endDateStr := req.EndDate.Format("2006-01-02")
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	req.StartDate = startDate
	req.EndDate = endDate
	booking, err := s.store.UpdateBooking(ctx, req)
	if err != nil {
		resp.Error = err.Error()
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, resp)
			return
		}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Booking = booking
	ctx.JSON(http.StatusOK, resp)
}

type deleteBookingRequest struct {
	ID int64 `json:"ID"`
}

func (s *Server) deleteBooking(ctx *gin.Context) {
	var req deleteBookingRequest
	var resp sqlc.DeleteBookingResponse
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	err := s.store.DeleteBooking(ctx, req.ID)
	if err != nil {
		resp.Error = err.Error()
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, resp)
		}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (s *Server) searchBookings(ctx *gin.Context) {
	query, err := url.QueryUnescape(ctx.Query("query"))
	var resp sqlc.BookingsResponse
	if err != nil {
		resp.Error = "invalid search query"
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	bookings, err := s.store.SearchBookings(ctx, sql.NullString{
		String: query,
		Valid:  true,
	})
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	// If the caller is unauthorized user, then details is empty
	if _, _, err = checkUserAuthorized(ctx, s.tokenMaker); err != nil {
		var anonymousBookings = make([]sqlc.Booking, 0, 1)
		for _, booking := range bookings {
			booking.Details = ""
			booking.TotalPrice = 0
			booking.CustomerName = ""
			anonymousBookings = append(anonymousBookings, booking)
			fmt.Printf("%+v\n", booking)
		}
		resp.Bookings = anonymousBookings
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.Bookings = bookings
	ctx.JSON(http.StatusOK, resp)
}
