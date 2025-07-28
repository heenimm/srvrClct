package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"serverCalc/handlers"
	"serverCalc/storage"

	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	absPath, err := filepath.Abs("/Users/halimlka/GolandProjects/serverCalc/calc.db")
	if err != nil {
		t.Fatalf("failed to resolve path: %v", err)
	}

	db, err := storage.InitDB(absPath)
	if err != nil {
		t.Log(err)
	}
	//_ = storage.ApplySchema(db, "../db/schema.sql")

	handler := handlers.RegisterHandler(db)
	reqBody := map[string]string{
		"login":    "stuser",
		"password": "secret123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	//if res.StatusCode != http.StatusCreated {
	//	body, _ := io.ReadAll(res.Body)
	//	t.Fatalf("expected status 201 Created, got %d: %s", res.StatusCode, string(body))
	//}
}
