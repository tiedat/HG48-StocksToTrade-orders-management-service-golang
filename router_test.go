package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	s, err := newServer()
	require.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			http.HandlerFunc(s.ping).ServeHTTP(tc.rec, tc.req)
			assert.Equal(t, tc.expectedBody, tc.rec.Body.String())
		})
	}
}
