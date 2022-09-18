package sqlc

import (
	"context"
	"database/sql"
	"github.com/mearaj/bhagad-house-booking/backend/db/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateBookingAndCustomer(t *testing.T) {
	store := NewStore(testDB)
	customerParams := CreateCustomerParams{
		Name:    util.RandomCustomerName(),
		Address: util.RandomCustomerAddr(),
		Phone:   util.RandomPhone(),
		Email:   util.RandomEmail(),
	}
	bookingParams := CreateBookingParams{
		StartDate: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		EndDate: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Rate: sql.NullFloat64{
			Float64: 1,
			Valid:   true,
		},
		RateTimeUnit: RateTimeUnitsDay,
	}

	n := 5
	errs := make(chan error)
	results := make(chan CreateBookingAndCustomerResult)
	for i := 0; i < 5; i++ {
		go func() {
			params := CreateBookingAndCustomerParams{
				CreateBookingParams:  bookingParams,
				CreateCustomerParams: customerParams,
			}
			result, err := store.CreateBookingAndCustomer(context.Background(), params)
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)
		require.NotZero(t, result.CustomerID)
		require.Equal(t, result.CustomerID.Int64, result.Customer.ID)
		_, err = store.GetBooking(context.Background(), result.Booking.ID)
		require.NoError(t, err)
		_, err = store.GetCustomer(context.Background(), result.Customer.ID)
		require.NoError(t, err)
	}
}
