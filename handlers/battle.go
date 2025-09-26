package handlers

import (
    "encoding/json"
    "net/http"
    "rpg-gx4-backend/battle"
    "rpg-gx4-backend/models"
    
    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Em produção, defina regras de origem adequadas
    },
}

func (h *AuthHandlers) CreateBattle(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)
    
    var req struct {
        OpponentID  string   `json:"opponent_id"`
        Characters []string `json:"characters"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    battle := h.battleManager.CreateBattle(userID, req.OpponentID, req.Characters)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(battle)
}

func (h *AuthHandlers) BattleWebsocket(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)
    battleID := mux.Vars(r)["id"]
    
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
        return
    }
    defer conn.Close()
    
    // Subscrever para atualizações de batalha
    battleChan := h.battleManager.Subscribe(battleID, userID)
    defer h.battleManager.Unsubscribe(battleID, userID)
    
    // Goroutine para enviar atualizações ao cliente
    go func() {
        for battle := range battleChan {
            if err := conn.WriteJSON(battle); err != nil {
                return
            }
        }
    }()
    
    // Ler ações do cliente
    for {
        var action models.BattleAction
        if err := conn.ReadJSON(&action); err != nil {
            return
        }
        
        action.BattleID = battleID
        action.PlayerID = userID
        h.battleManager.SubmitAction(action)
    }
}

func (h *AuthHandlers) SubmitBattleAction(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)
    battleID := mux.Vars(r)["id"]
    
    var action models.BattleAction
    if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    action.BattleID = battleID
    action.PlayerID = userID
    
    h.battleManager.SubmitAction(action)
    w.WriteHeader(http.StatusAccepted)
}