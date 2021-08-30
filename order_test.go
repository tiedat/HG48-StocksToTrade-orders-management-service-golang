package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrdersHandler(t *testing.T) {
	tests := []struct {
		name           string
		rec            *httptest.ResponseRecorder
		req            *http.Request
		expectedStatus int
	}{
		{
			"Success",
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "/orders", nil),
			200,
		},
	}

	s, err := newServer(withDB(db))
	require.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s.hs.Handler.ServeHTTP(tc.rec, tc.req)
			assert.Equal(t, tc.expectedStatus, tc.rec.Code)
			t.Log(tc.rec.Body)
		})
	}
}

func TestGetOrderDetail(t *testing.T) {
	tests := []struct {
		name           string
		rec            *httptest.ResponseRecorder
		req            *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			"Invalid URL",
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "/orders/", nil),
			404,
			"404 page not found\n",
		},
		{
			"not-existent email",
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "/orders/abc", nil),
			404,
			"\"email not found\"\n",
		},
		{
			"Success",
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "/orders/foo@example.com", nil),
			200,
			"{\"id\":null,\"email\":\"foo@example.com\",\"product_id\":null,\"create_at\":null,\"updated_at\":null}\n",
		},
	}

	s, err := newServer(withDB(db))
	require.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s.hs.Handler.ServeHTTP(tc.rec, tc.req)
			assert.Equal(t, tc.expectedStatus, tc.rec.Code)
			assert.Equal(t, tc.expectedBody, tc.rec.Body.String())
		})
	}
}

func TestCancelSubscription(t *testing.T) {
	tableTest := []struct {
		name          string
		email         string
		expectedError error
	}{
		{
			"not-existent email",
			"",
			errors.New(fmt.Sprintf("user not found for %v", "")),
		},
		{
			"successfully insert entitlement log",
			"foo@example.com",
			nil,
		},
		{
			"successfully update nyse_entries",
			"bar@example.com",
			nil,
		},
		{
			"successfully update referrals",
			"fizz@example.com",
			nil,
		},
	}

	s, err := newServer(withDB(db))
	require.NoError(t, err)

	for _, tc := range tableTest {
		t.Run(tc.name, func(t *testing.T) {
			err := s.cancelSubscription(tc.email)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestReactiveSubscription(t *testing.T) {
	tableTest := []struct {
		name          string
		email         string
		exceptedError error
	}{
		{
			"not-existent email",
			"",
			errors.New(fmt.Sprintf("user not found for %v", "")),
		},
		{
			"successfully update users",
			"foo@example.com",
			nil,
		},
	}

	s, err := newServer(withDB(db))
	require.NoError(t, err)

	for _, tc := range tableTest {
		t.Run(tc.name, func(t *testing.T) {
			err := s.reactiveSubscription(tc.email)
			assert.Equal(t, tc.exceptedError, err)
		})
	}
}

func TestGetRecurlySubscription(t *testing.T) {
	tableTest := []struct {
		name           string
		subscriptionId string
		exceptedData   []*RecurlySubscription
		exceptedError  error
	}{
		{
			"not-existent subscription id",
			"",
			nil,
			nil,
		},
		{
			"successfully get recurly subscription",
			"4f5d7e1db40b80b1aa319a418a8fe9de",
			[]*RecurlySubscription{
				{
					ID: 1, UserId: 3805, RecurlySubscriptionId: "4f5d7e1db40b80b1aa319a418a8fe9de",
				},
			},
			nil,
		},
	}

	s, err := newServer(withDB(db))
	require.NoError(t, err)

	for _, tc := range tableTest {
		t.Run(tc.name, func(t *testing.T) {
			data, err := s.getRecurlySubscription(tc.subscriptionId)
			assert.Equal(t, tc.exceptedData, data)
			assert.Equal(t, tc.exceptedError, err)
		})

	}
}
