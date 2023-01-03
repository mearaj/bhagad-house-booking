package sqlc

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const overlappingRangeError = "pq: conflicting key value violates exclusion constraint \"nooverlappingtimeranges\""

func TestCreateBooking(t *testing.T) {
	//https://stackoverflow.com/questions/60433870/saving-time-time-in-golang-to-postgres-timestamp-with-time-zone-field
	arg := CreateBookingParams{
		StartDate: time.Now().UTC().AddDate(0, 0, 1).Round(time.Microsecond),
		EndDate:   time.Now().UTC().AddDate(0, 0, 14).Round(time.Microsecond),
		Details:   "",
	}
	bkg, err := testQueries.CreateBooking(context.Background(), arg)
	for err != nil && err.Error() == overlappingRangeError {
		arg.StartDate = arg.StartDate.Add(time.Hour * 24)
		arg.EndDate = arg.EndDate.Add(time.Hour * 24)
		bkg, err = testQueries.CreateBooking(context.Background(), arg)
	}
	require.NoError(t, err)
	require.NotEmpty(t, bkg)
	require.True(t, bkg.StartDate.Equal(arg.StartDate))
	require.True(t, bkg.EndDate.Equal(arg.EndDate))
	require.True(t, bkg.Details == arg.Details)
}
