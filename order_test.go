package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrder(t *testing.T) {
	tests := []struct {
		name           string
		rec            *httptest.ResponseRecorder
		req            *http.Request
		expectedStatus int
	}{
		{
			"ok",
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "/order/getAll", nil),
			200,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			http.HandlerFunc(getAllOrderHandler).ServeHTTP(tc.rec, tc.req)
			fmt.Print(tc.rec.Body)
		})
	}
}
