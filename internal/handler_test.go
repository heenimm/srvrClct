package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"serverCalc/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddExpressionHandler(t *testing.T) {
	requestBody, _ := json.Marshal(pkg.CalculationRequest{Expression: "2+2"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddExpressionHandler)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code, "должен возвращать 201 Created")
	var response pkg.CalculationResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err, "ответ должен быть корректным JSON")
	assert.NotEmpty(t, response.ID, "ID не должен быть пустым")
}

func TestAddExpressionHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddExpressionHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "должен возвращать 405 Method Not Allowed")
}

func TestAddExpressionHandler_InvalidRequestBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer([]byte("{invalid_json}")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddExpressionHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code, "должен возвращать 422 Unprocessable Entity")
}

func TestGetExpressionsHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/expressions", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetExpressionsHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "должен возвращать 405 Method Not Allowed")
}

func TestGetExpressionsHandler_StoreError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/expressions", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetExpressionsHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "должен возвращать 500 Internal Server Error")
}
