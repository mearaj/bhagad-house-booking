// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mearaj/bhagad-house-booking/common/db/sqlc (interfaces: Store)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	sqlc "github.com/mearaj/bhagad-house-booking/common/db/sqlc"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateBooking mocks base method.
func (m *MockStore) CreateBooking(arg0 context.Context, arg1 sqlc.CreateBookingParams) (sqlc.Booking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBooking", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Booking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBooking indicates an expected call of CreateBooking.
func (mr *MockStoreMockRecorder) CreateBooking(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBooking", reflect.TypeOf((*MockStore)(nil).CreateBooking), arg0, arg1)
}

// CreateBookingAndCustomer mocks base method.
func (m *MockStore) CreateBookingAndCustomer(arg0 context.Context, arg1 sqlc.CreateBookingAndCustomerParams) (sqlc.CreateBookingAndCustomerResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBookingAndCustomer", arg0, arg1)
	ret0, _ := ret[0].(sqlc.CreateBookingAndCustomerResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBookingAndCustomer indicates an expected call of CreateBookingAndCustomer.
func (mr *MockStoreMockRecorder) CreateBookingAndCustomer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBookingAndCustomer", reflect.TypeOf((*MockStore)(nil).CreateBookingAndCustomer), arg0, arg1)
}

// CreateCustomer mocks base method.
func (m *MockStore) CreateCustomer(arg0 context.Context, arg1 sqlc.CreateCustomerParams) (sqlc.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCustomer", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCustomer indicates an expected call of CreateCustomer.
func (mr *MockStoreMockRecorder) CreateCustomer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCustomer", reflect.TypeOf((*MockStore)(nil).CreateCustomer), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 sqlc.CreateUserParams) (sqlc.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(sqlc.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// DeleteBooking mocks base method.
func (m *MockStore) DeleteBooking(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBooking", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBooking indicates an expected call of DeleteBooking.
func (mr *MockStoreMockRecorder) DeleteBooking(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBooking", reflect.TypeOf((*MockStore)(nil).DeleteBooking), arg0, arg1)
}

// DeleteCustomer mocks base method.
func (m *MockStore) DeleteCustomer(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCustomer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCustomer indicates an expected call of DeleteCustomer.
func (mr *MockStoreMockRecorder) DeleteCustomer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCustomer", reflect.TypeOf((*MockStore)(nil).DeleteCustomer), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockStore) DeleteUser(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStoreMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), arg0, arg1)
}

// GetBooking mocks base method.
func (m *MockStore) GetBooking(arg0 context.Context, arg1 int64) (sqlc.Booking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBooking", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Booking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBooking indicates an expected call of GetBooking.
func (mr *MockStoreMockRecorder) GetBooking(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBooking", reflect.TypeOf((*MockStore)(nil).GetBooking), arg0, arg1)
}

// GetCustomer mocks base method.
func (m *MockStore) GetCustomer(arg0 context.Context, arg1 int64) (sqlc.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCustomer", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCustomer indicates an expected call of GetCustomer.
func (mr *MockStoreMockRecorder) GetCustomer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomer", reflect.TypeOf((*MockStore)(nil).GetCustomer), arg0, arg1)
}

// GetUserByEmail mocks base method.
func (m *MockStore) GetUserByEmail(arg0 context.Context, arg1 string) (sqlc.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(sqlc.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockStoreMockRecorder) GetUserByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockStore)(nil).GetUserByEmail), arg0, arg1)
}

// GetUserByID mocks base method.
func (m *MockStore) GetUserByID(arg0 context.Context, arg1 int64) (sqlc.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0, arg1)
	ret0, _ := ret[0].(sqlc.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockStoreMockRecorder) GetUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockStore)(nil).GetUserByID), arg0, arg1)
}

// ListBookings mocks base method.
func (m *MockStore) ListBookings(arg0 context.Context, arg1 sqlc.ListBookingsParams) ([]sqlc.Booking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBookings", arg0, arg1)
	ret0, _ := ret[0].([]sqlc.Booking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBookings indicates an expected call of ListBookings.
func (mr *MockStoreMockRecorder) ListBookings(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBookings", reflect.TypeOf((*MockStore)(nil).ListBookings), arg0, arg1)
}

// ListCustomers mocks base method.
func (m *MockStore) ListCustomers(arg0 context.Context, arg1 sqlc.ListCustomersParams) ([]sqlc.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCustomers", arg0, arg1)
	ret0, _ := ret[0].([]sqlc.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCustomers indicates an expected call of ListCustomers.
func (mr *MockStoreMockRecorder) ListCustomers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCustomers", reflect.TypeOf((*MockStore)(nil).ListCustomers), arg0, arg1)
}

// ListUsers mocks base method.
func (m *MockStore) ListUsers(arg0 context.Context, arg1 sqlc.ListUsersParams) ([]sqlc.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", arg0, arg1)
	ret0, _ := ret[0].([]sqlc.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockStoreMockRecorder) ListUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockStore)(nil).ListUsers), arg0, arg1)
}

// UpdateBooking mocks base method.
func (m *MockStore) UpdateBooking(arg0 context.Context, arg1 sqlc.UpdateBookingParams) (sqlc.Booking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBooking", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Booking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBooking indicates an expected call of UpdateBooking.
func (mr *MockStoreMockRecorder) UpdateBooking(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBooking", reflect.TypeOf((*MockStore)(nil).UpdateBooking), arg0, arg1)
}

// UpdateCustomer mocks base method.
func (m *MockStore) UpdateCustomer(arg0 context.Context, arg1 sqlc.UpdateCustomerParams) (sqlc.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCustomer", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCustomer indicates an expected call of UpdateCustomer.
func (mr *MockStoreMockRecorder) UpdateCustomer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCustomer", reflect.TypeOf((*MockStore)(nil).UpdateCustomer), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockStore) UpdateUser(arg0 context.Context, arg1 sqlc.UpdateUserParams) (sqlc.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(sqlc.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockStoreMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockStore)(nil).UpdateUser), arg0, arg1)
}
