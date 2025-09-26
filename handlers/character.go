package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"rpg-gx4-backend/models"
)

func (h *AuthHandlers) CreateCharacter(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var char models.Character
	if err := json.NewDecoder(r.Body).Decode(&char); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validar dados básicos
	if char.Name == "" || char.Class == 0 {
		http.Error(w, "Name and class are required", http.StatusBadRequest)
		return
	}

	// Inserir no banco de dados
	err := h.db.QueryRow(`
		INSERT INTO characters (
			user_id, name, level, xp, win, lose, class, subclass, 
			subclass_traits, spells, is_summoner, race, gender, 
			color, talents, weapons, attack_effect
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id, created_at
	`,
		userID, char.Name, char.Level, char.XP, char.Win, char.Lose, char.Class, char.Subclass,
		char.SubclassTraits, char.Spells, char.IsSummoner, char.Race, char.Gender,
		char.Color, char.Talents, char.Weapons, char.AttackEffect,
	).Scan(&char.ID, &char.CreatedAt)

	if err != nil {
		http.Error(w, "Failed to create character", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(char)
}

func (h *AuthHandlers) GetUserCharacters(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	rows, err := h.db.Query(`
		SELECT id, name, level, xp, win, lose, class, subclass, 
		subclass_traits, spells, is_summoner, race, gender, 
		color, talents, weapons, attack_effect, created_at
		FROM characters WHERE user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var characters []models.Character
	for rows.Next() {
		var char models.Character
		err := rows.Scan(
			&char.ID, &char.Name, &char.Level, &char.XP, &char.Win, &char.Lose, &char.Class, &char.Subclass,
			&char.SubclassTraits, &char.Spells, &char.IsSummoner, &char.Race, &char.Gender,
			&char.Color, &char.Talents, &char.Weapons, &char.AttackEffect, &char.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		characters = append(characters, char)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}

func (h *AuthHandlers) GetCharacter(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	charID := mux.Vars(r)["id"]

	var char models.Character
	err := h.db.QueryRow(`
		SELECT id, user_id, name, level, xp, win, lose, class, subclass, 
		subclass_traits, spells, is_summoner, race, gender, 
		color, talents, weapons, attack_effect, created_at
		FROM characters WHERE id = $1 AND user_id = $2
	`, charID, userID).Scan(
		&char.ID, &char.UserID, &char.Name, &char.Level, &char.XP, &char.Win, &char.Lose, &char.Class, &char.Subclass,
		&char.SubclassTraits, &char.Spells, &char.IsSummoner, &char.Race, &char.Gender,
		&char.Color, &char.Talents, &char.Weapons, &char.AttackEffect, &char.CreatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(char)
}

func (h *AuthHandlers) UpdateCharacter(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	charID := mux.Vars(r)["id"]

	var char models.Character
	if err := json.NewDecoder(r.Body).Decode(&char); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Atualizar no banco de dados
	_, err := h.db.Exec(`
		UPDATE characters SET
			name = $1,
			level = $2,
			xp = $3,
			win = $4,
			lose = $5,
			class = $6,
			subclass = $7,
			subclass_traits = $8,
			spells = $9,
			is_summoner = $10,
			race = $11,
			gender = $12,
			color = $13,
			talents = $14,
			weapons = $15,
			attack_effect = $16
		WHERE id = $17 AND user_id = $18
	`,
		char.Name, char.Level, char.XP, char.Win, char.Lose, char.Class, char.Subclass,
		char.SubclassTraits, char.Spells, char.IsSummoner, char.Race, char.Gender,
		char.Color, char.Talents, char.Weapons, char.AttackEffect,
		charID, userID,
	)

	if err != nil {
		http.Error(w, "Failed to update character", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandlers) DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	charID := mux.Vars(r)["id"]

	_, err := h.db.Exec("DELETE FROM characters WHERE id = $1 AND user_id = $2", charID, userID)
	if err != nil {
		http.Error(w, "Failed to delete character", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}