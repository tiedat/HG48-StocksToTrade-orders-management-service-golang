package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name         string
		rec          *httptest.ResponseRecorder
		req          *http.Request
		expectedBody string
	}{
		{
			"ok",
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "/ping", nil),
			"pong",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			http.HandlerFunc(ping).ServeHTTP(tc.rec, tc.req)

			assert.Equal(t, tc.expectedBody, tc.rec.Body.String())
		})
	}
}
