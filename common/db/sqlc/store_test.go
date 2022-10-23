package sqlc

import (
	"context"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateBookingAndCustomer(t *testing.T) {
	store := NewStore(testDB)
	customerParams := CreateCustomerParams{
		Name:    utils.RandomName(),
		Address: utils.RandomAddress(),
		Phone:   utils.RandomPhone(),
		Email:   utils.RandomEmail(),
	}
	bookingParams := CreateBookingParams{
		StartDate:    time.Now(),
		EndDate:      time.Now(),
		Rate:         1,
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
		require.Equal(t, result.CustomerID, result.Customer.ID)
		_, err = store.GetBooking(context.Background(), result.Booking.ID)
		require.NoError(t, err)
		_, err = store.GetCustomer(context.Background(), result.Customer.ID)
		require.NoError(t, err)
	}
}
