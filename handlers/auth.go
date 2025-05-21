package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"rpg-backend/models"
)

type AuthHandlers struct {
	db *sql.DB
}

func NewHandlers(db *sql.DB) *AuthHandlers {
	return &AuthHandlers{db: db}
}

type GoogleAuthRequest struct {
	Token string `json:"token"`
}

func (h *AuthHandlers) HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	var req GoogleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Aqui você validaria o token do Google e extrairia o google_sub
	// Por simplicidade, vou assumir que o token é válido e o google_sub está no token
	googleSub := req.Token // Na prática, você extrairia isso do token JWT

	// Verificar se o usuário já existe
	var user models.User
	err := h.db.QueryRow("SELECT id, google_sub, created_at FROM users WHERE google_sub = $1", googleSub).
		Scan(&user.ID, &user.GoogleSub, &user.CreatedAt)

	if err == sql.ErrNoRows {
		// Criar novo usuário
		err = h.db.QueryRow(
			"INSERT INTO users (google_sub) VALUES ($1) RETURNING id, created_at",
			googleSub,
		).Scan(&user.ID, &user.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Retornar informações do usuário (sem o google_sub)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         user.ID,
		"created_at": user.CreatedAt,
	})
}

func (h *AuthHandlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extrair token do header Authorization
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validar token e obter google_sub
		googleSub := token // Na prática, você validaria o token JWT aqui

		// Verificar se o usuário existe
		var userID string
		err := h.db.QueryRow("SELECT id FROM users WHERE google_sub = $1", googleSub).Scan(&userID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Adicionar userID ao contexto da requisição
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}