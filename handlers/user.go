package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"rpg-gx4-backend/models"
)

func GetUserByGoogleSub(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		googleSub := r.URL.Query().Get("google_sub")
		if googleSub == "" {
			http.Error(w, "Google sub is required", http.StatusBadRequest)
			return
		}

		var user models.User
		err := db.QueryRow("SELECT id, google_sub, created_at FROM users WHERE google_sub = $1", googleSub).
			Scan(&user.ID, &user.GoogleSub, &user.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.GoogleSub == "" {
			http.Error(w, "Google sub is required", http.StatusBadRequest)
			return
		}

		var id string
		err := db.QueryRow(
			"INSERT INTO users (google_sub) VALUES ($1) RETURNING id",
			user.GoogleSub,
		).Scan(&id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.ID = id
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}