package sqlc

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateBooking(t *testing.T) {
	cust := createRandomCustomer(t)
	arg := CreateBookingParams{
		StartDate:    time.Time{},
		EndDate:      time.Time{},
		CustomerID:   cust.ID,
		Rate:         1,
		RateTimeUnit: RateTimeUnitsDay,
	}
	bkg, err := testQueries.CreateBooking(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, bkg)
	require.Equal(t, bkg.StartDate, arg.StartDate)
	require.Equal(t, bkg.EndDate, arg.EndDate)
	require.Equal(t, bkg.CustomerID, arg.CustomerID)
	require.Equal(t, bkg.Rate, arg.Rate)
	require.Equal(t, bkg.Rate, arg.Rate)
	require.Equal(t, bkg.RateTimeUnit, arg.RateTimeUnit)
}
