package internal_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"serverCalc/internal"
	"testing"
)

func TestCalculateHandler_Success(t *testing.T) {

	// arrange //
	////////////

	// Создаем запрос с корректным телом
	reqBody := `{"expression": "2+2*2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для захвата ответа
	rec := httptest.NewRecorder()

	// act //
	/////////

	// Вызываем реальный обработчик
	handler := http.HandlerFunc(internal.CalculateHandler)
	handler.ServeHTTP(rec, req)

	// assert //
	////////////

	// Проверяем ответ
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var response struct {
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Result != 6 {
		t.Errorf("expected result 6, got %f", response.Result)
	}
}

func TestCalculateHandler_UnprocessableEntity(t *testing.T) {
	// arrange //
	////////////

	// Некорректное тело запроса
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

	// assert //
	////////////

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", res.StatusCode)
	}

	var response struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error != "Expression is not valid" {
		t.Errorf("expected error 'Expression is not valid', got '%s'", response.Error)
	}
}

func TestCalculateHandler_BadRequest(t *testing.T) {
	// arrange //
	////////////

	// Деление на ноль
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

	// assert //
	////////////

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.StatusCode)
	}

	var response struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error != "Internal server error" {
		t.Errorf("expected error 'Internal server error', got '%s'", response.Error)
	}
}
