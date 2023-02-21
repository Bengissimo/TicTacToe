package game

import (
	"github.com/google/uuid"
)

const (
	STATUS_RUNNING = "RUNNING"
	STATUS_X_WON   = "X_WON"
	STATUS_O_WON   = "O_WON"
	STATUS_DRAW    = "DRAW"
)

type Game struct {
	ID           uuid.UUID `json:"id"`
	Board        string    `json:"board" binding:"required,len=9"`
	Status       string    `json:"status"`
	ServerSymbol byte
}
