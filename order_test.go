package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
			"ok",
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "/orders", nil),
			200,
		},
	}

	s, err := newServer(withDB(db))
	require.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			http.HandlerFunc(s.ordersHandler).ServeHTTP(tc.rec, tc.req)
			t.Log(tc.rec.Body)
		})
	}
}
