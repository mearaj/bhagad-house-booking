package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateIndianPhoneNumber(t *testing.T) {
	type TestNumber struct {
		Number string
		Valid  bool
	}
	numbers := []TestNumber{
		{Number: "+911234567890", Valid: true},
		{Number: "+91234567890", Valid: false},
		{Number: "11234567890", Valid: false},
		{Number: "01234567890", Valid: true},
		{Number: "234567890", Valid: false},
		{Number: "911234567890", Valid: true},
		{Number: "+91123456789", Valid: false},
	}
	for _, num := range numbers {
		res := ValidateIndianPhoneNumber(num.Number)
		require.Equal(t, res, num.Valid)
	}
}
