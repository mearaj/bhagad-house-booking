// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mearaj/bhagad-house-booking/common/db/sqlc (interfaces: Store)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	sql "database/sql"
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

// GetConflictingBookings mocks base method.
func (m *MockStore) GetConflictingBookings(arg0 context.Context, arg1 sqlc.GetConflictingBookingsParams) ([]sqlc.Booking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConflictingBookings", arg0, arg1)
	ret0, _ := ret[0].([]sqlc.Booking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConflictingBookings indicates an expected call of GetConflictingBookings.
func (mr *MockStoreMockRecorder) GetConflictingBookings(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConflictingBookings", reflect.TypeOf((*MockStore)(nil).GetConflictingBookings), arg0, arg1)
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

// SearchBookings mocks base method.
func (m *MockStore) SearchBookings(arg0 context.Context, arg1 sql.NullString) ([]sqlc.Booking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchBookings", arg0, arg1)
	ret0, _ := ret[0].([]sqlc.Booking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchBookings indicates an expected call of SearchBookings.
func (mr *MockStoreMockRecorder) SearchBookings(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchBookings", reflect.TypeOf((*MockStore)(nil).SearchBookings), arg0, arg1)
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
