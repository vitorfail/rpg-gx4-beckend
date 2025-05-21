package models

import "time"

type Battle struct {
    ID           string    `json:"id"`
    Player1ID    string    `json:"player1_id"`
    Player2ID    string    `json:"player2_id"`
    Characters   []string  `json:"characters"` // IDs dos personagens participantes
    CurrentTurn  int       `json:"current_turn"`
    Status       string    `json:"status"`     // "waiting", "in_progress", "finished"
    Winner       *string   `json:"winner"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type BattleAction struct {
    BattleID   string `json:"battle_id"`
    PlayerID   string `json:"player_id"`
    CharacterID string `json:"character_id"`
    ActionType string `json:"action_type"` // "attack", "defend", "use_item", "cast_spell"
    TargetID   string `json:"target_id"`   // ID do alvo (personagem ou item)
    Value      int    `json:"value"`       // Valor do dano/cura/etc
}