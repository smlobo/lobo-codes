package test

import (
	"client-server/internal"
	"client-server/internal/handler"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	req, err := http.NewRequest("GET", internal.AsiaPath, nil)
	require.NoError(t, err, "Failed creating an HTTP GET request to:",
		internal.AsiaPath)

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record
	// the response.
	responseRecorder := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(handler.ReturnHandler("asia"))

	// Call the handler ServeHTTP method directly and pass in our Request and
	// ResponseRecorder.
	handlerFunc.ServeHTTP(responseRecorder, req)

	// Check the status code is what we expect.
	t.Log("Asia handler returned status code:", responseRecorder.Code)
	require.Equal(t, http.StatusOK, responseRecorder.Code)
}