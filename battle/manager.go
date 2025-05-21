package battle

import (
    "sync"
    "time"
    "rpg-backend/models"
)

type BattleManager struct {
    battles      map[string]*models.Battle
    actions      chan models.BattleAction
    subscriptions map[string]map[string]chan models.Battle // battleID -> userID -> channel
    mutex        sync.RWMutex
}

func NewBattleManager() *BattleManager {
    return &BattleManager{
        battles:      make(map[string]*models.Battle),
        actions:      make(chan models.BattleAction, 100),
        subscriptions: make(map[string]map[string]chan models.Battle),
    }
}

func (bm *BattleManager) Start() {
    go bm.processActions()
}

func (bm *BattleManager) processActions() {
    for action := range bm.actions {
        bm.mutex.Lock()
        
        battle, exists := bm.battles[action.BattleID]
        if !exists || battle.Status != "in_progress" {
            bm.mutex.Unlock()
            continue
        }

        // Processar ação (simplificado)
        // Aqui você implementaria a lógica real do combate
        switch action.ActionType {
        case "attack":
            // Lógica de ataque
        case "defend":
            // Lógica de defesa
        // etc...
        }

        battle.UpdatedAt = time.Now()
        
        // Notificar todos os subscribers
        for _, sub := range bm.subscriptions[action.BattleID] {
            select {
            case sub <- *battle:
            default:
                // Canal cheio, pode querer lidar com isso
            }
        }
        
        bm.mutex.Unlock()
    }
}

func (bm *BattleManager) CreateBattle(player1ID, player2ID string, characters []string) *models.Battle {
    battle := &models.Battle{
        ID:          generateBattleID(),
        Player1ID:   player1ID,
        Player2ID:   player2ID,
        Characters:  characters,
        Status:      "waiting",
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    bm.mutex.Lock()
    bm.battles[battle.ID] = battle
    bm.subscriptions[battle.ID] = make(map[string]chan models.Battle)
    bm.mutex.Unlock()
    
    return battle
}

func (bm *BattleManager) Subscribe(battleID, userID string) chan models.Battle {
    ch := make(chan models.Battle, 10)
    
    bm.mutex.Lock()
    bm.subscriptions[battleID][userID] = ch
    bm.mutex.Unlock()
    
    return ch
}

func (bm *BattleManager) Unsubscribe(battleID, userID string) {
    bm.mutex.Lock()
    if subs, ok := bm.subscriptions[battleID]; ok {
        if ch, ok := subs[userID]; ok {
            close(ch)
            delete(subs, userID)
        }
    }
    bm.mutex.Unlock()
}

func (bm *BattleManager) SubmitAction(action models.BattleAction) {
    bm.actions <- action
}

func generateBattleID() string {
    // Implementar geração de ID (pode usar UUID como você já faz)
    return "battle_" + time.Now().Format("20060102150405")
}