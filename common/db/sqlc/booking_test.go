package sqlc

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateBooking(t *testing.T) {
	cust := createRandomCustomer(t)
	arg := CreateBookingParams{
		StartDate: sql.NullTime{
			Valid: true,
		},
		EndDate: sql.NullTime{
			Valid: true,
		},
		CustomerID: sql.NullInt64{
			Int64: cust.ID,
			Valid: true,
		},
		Rate: sql.NullFloat64{
			Float64: 1,
			Valid:   true,
		},
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
