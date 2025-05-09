package handlers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegisterHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if db == nil {
			http.Error(w, "Database not initialized", http.StatusInternalServerError)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// Хешируем пароль
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "error hashing password", http.StatusInternalServerError)
			return
		}

		// Пишем в БД
		_, err = db.Exec("INSERT INTO users (login, password_hash) VALUES (?, ?)", req.Login, string(hash))
		if err != nil {
			http.Error(w, "user already exists or db error", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("user registered"))
	}
}
