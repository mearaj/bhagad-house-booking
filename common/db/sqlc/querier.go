// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package sqlc

import (
	"context"
)

type Querier interface {
	CreateBooking(ctx context.Context, arg CreateBookingParams) (Booking, error)
	CreateCustomer(ctx context.Context, arg CreateCustomerParams) (Customer, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteBooking(ctx context.Context, id int64) error
	DeleteCustomer(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetBooking(ctx context.Context, id int64) (Booking, error)
	GetCustomer(ctx context.Context, id int64) (Customer, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	ListBookings(ctx context.Context, arg ListBookingsParams) ([]Booking, error)
	ListCustomers(ctx context.Context, arg ListCustomersParams) ([]Customer, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateBooking(ctx context.Context, arg UpdateBookingParams) (Booking, error)
	UpdateCustomer(ctx context.Context, arg UpdateCustomerParams) (Customer, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)