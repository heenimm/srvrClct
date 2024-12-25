package internal_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"serverCalc/internal"
	"testing"
)

func TestCalculateHandler_Success(t *testing.T) {

	// arrange //
	////////////

	reqBody := `{"expression": "2+2*2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// act //
	/////////

	internal.CalculateHandler(rec, req)

	// assert //
	////////////

	res := rec.Result()
	defer res.Body.Close()

	var response struct {
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	require.Equal(t, res.StatusCode, http.StatusOK)
	require.Equal(t, response.Result, float64(6))
}

func TestCalculateHandler_UnprocessableEntity(t *testing.T) {
	// arrange //
	////////////

	reqBody := `{"Expression": "2+2*abc"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// act //
	/////////

	handler := http.HandlerFunc(internal.CalculateHandler)
	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	var response struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// assert //
	////////////

	require.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
	require.Equal(t, response.Error, "Expression is not valid")
}

func TestCalculateHandler_BadRequest(t *testing.T) {
	// arrange //
	////////////

	reqBody := `{"expression": "test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// act //
	/////////

	handler := http.HandlerFunc(internal.CalculateHandler)
	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	var response struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// assert //
	////////////

	require.Equal(t, res.StatusCode, http.StatusInternalServerError)
	require.Equal(t, response.Error, "Internal server error")
}
