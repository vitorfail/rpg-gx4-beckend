package models

import (
	"time"
)

type Character struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Name          string    `json:"name"`
	Level         int       `json:"level"`
	XP            int       `json:"xp"`
	Win           int       `json:"win"`
	Lose          int       `json:"lose"`
	Class         int       `json:"class"`
	Subclass      *int      `json:"subclass,omitempty"`
	SubclassTraits []int    `json:"subclass_traits,omitempty"`
	Spells        []int     `json:"spells,omitempty"`
	IsSummoner    bool      `json:"is_summoner"`
	Race          *int      `json:"race,omitempty"`
	Gender        *int      `json:"gender,omitempty"`
	Color         []float64 `json:"color,omitempty"`
	Talents       []int     `json:"talents,omitempty"`
	Weapons       []int     `json:"weapons,omitempty"`
	AttackEffect  *int      `json:"attack_effect,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}