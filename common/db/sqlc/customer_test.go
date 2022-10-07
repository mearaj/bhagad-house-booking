package sqlc

import (
	"context"
	"database/sql"
	"github.com/mearaj/bhagad-house-booking/common/db/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomCustomer(t *testing.T) Customer {
	arg := CreateCustomerParams{
		Name:    util.RandomCustomerName(),
		Address: util.RandomCustomerAddr(),
		Phone:   util.RandomPhone(),
		Email:   util.RandomEmail(),
	}
	customer, err := testQueries.CreateCustomer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, customer)
	require.Equal(t, customer.Name, arg.Name)
	require.Equal(t, customer.Address, arg.Address)
	require.Equal(t, customer.Phone, arg.Phone)
	require.Equal(t, customer.Email, arg.Email)
	return customer
}

func TestCreateCustomer(t *testing.T) {
	createRandomCustomer(t)
}

func TestGetCustomer(t *testing.T) {
	customer := createRandomCustomer(t)
	customer2, err := testQueries.GetCustomer(context.Background(), customer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, customer2)
	require.Equal(t, customer.ID, customer2.ID)
	require.Equal(t, customer.Name, customer2.Name)
	require.Equal(t, customer.Address, customer2.Address)
	require.Equal(t, customer.Phone, customer2.Phone)
	require.Equal(t, customer.Email, customer2.Email)
}

func TestUpdateCustomer(t *testing.T) {
	customer := createRandomCustomer(t)
	arg := UpdateCustomerParams{
		ID:      customer.ID,
		Name:    "Rahim",
		Address: "Bhagad house",
		Phone:   "123",
		Email:   "rahim@bhagad.com",
	}
	customer2, err := testQueries.UpdateCustomer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, customer2)
	require.Equal(t, customer.ID, customer2.ID)
	require.Equal(t, arg.Name, customer2.Name)
	require.Equal(t, arg.Address, customer2.Address)
	require.Equal(t, arg.Phone, customer2.Phone)
	require.Equal(t, arg.Email, customer2.Email)
}

func TestDeleteCustomer(t *testing.T) {
	customer := createRandomCustomer(t)
	err := testQueries.DeleteCustomer(context.Background(), customer.ID)
	require.NoError(t, err)

	customer2, err := testQueries.GetCustomer(context.Background(), customer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, customer2)
}

func TestListCustomers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCustomer(t)
	}
	arg := ListCustomersParams{
		Limit:  5,
		Offset: 5,
	}
	customers, err := testQueries.ListCustomers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, customers, 5)
	for _, customer := range customers {
		require.NotEmpty(t, customer)
	}
}
