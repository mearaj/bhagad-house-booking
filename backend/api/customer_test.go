package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/mearaj/bhagad-house-booking/common/db/mock"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/common/token"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func randomCustomer() sqlc.Customer {
	return sqlc.Customer{
		ID:        utils.RandomInt(1, 1000),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Name:      utils.RandomName(),
		Address:   utils.RandomAddress(),
		Phone:     utils.RandomPhone(),
		Email:     utils.RandomEmail(),
	}
}

func TestGetCustomerAPI(t *testing.T) {
	customer := randomCustomer()

	testCases := []struct {
		name          string
		customerID    int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "mearajbhagad@gmail.com", time.Minute)
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCustomer(t, recorder.Body, customer)
			},
		},
		{
			name:       "NotFound",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "mearajbhagad@gmail.com", time.Minute)
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(sqlc.Customer{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "mearajbhagad@gmail.com", time.Minute)
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(sqlc.Customer{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidID",
			customerID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "mearajbhagad@gmail.com", time.Minute)
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mock.NewMockStore(ctrl)
			tc.buildStubs(store)
			// start test server and request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/customers/%d", tc.customerID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchCustomer(t *testing.T, body *bytes.Buffer, customer sqlc.Customer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var gotCustomer sqlc.Customer
	err = json.Unmarshal(data, &gotCustomer)
	require.NoError(t, err)
	require.Equal(t, customer, gotCustomer)
}
