package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	CreateBookingAndCustomer(ctx context.Context, arg CreateBookingAndCustomerParams) (result CreateBookingAndCustomerResult, err error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}
	return tx.Commit()
}

type CreateBookingAndCustomerParams struct {
	CreateBookingParams
	CreateCustomerParams
}

type CreateBookingAndCustomerResult struct {
	Booking
	Customer
}

// CreateBookingAndCustomer
func (store *SQLStore) CreateBookingAndCustomer(ctx context.Context, arg CreateBookingAndCustomerParams) (result CreateBookingAndCustomerResult, err error) {
	err = store.execTx(ctx, func(q *Queries) (err error) {
		var customer Customer
		customer, err = q.CreateCustomer(ctx, arg.CreateCustomerParams)
		if err != nil {
			return err
		}
		result.Customer = customer
		arg.CreateBookingParams.CustomerID = customer.ID
		var booking Booking
		booking, err = q.CreateBooking(ctx, arg.CreateBookingParams)
		if err != nil {
			return err
		}
		result.Booking = booking
		return nil
	})
	return result, err
}
