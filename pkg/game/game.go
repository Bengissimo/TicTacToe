package game

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

const (
	STATUS_RUNNING = "RUNNING"
	STATUS_X_WON   = "X_WON"
	STATUS_O_WON   = "O_WON"
	STATUS_DRAW    = "DRAW"
	BOARD_LEN      = 9
	SYMBOL_X       = 'X'
	SYMBOL_O       = 'O'
	EMPTY          = '-'
)

type Game struct {
	ID              uuid.UUID `json:"id"`
	Board           string    `json:"board" binding:"required,len=9"`
	Status          string    `json:"status"`
	serverSymbol    byte
	clientSymbol    byte
	randomGenerator *rand.Rand
}

func (g *Game) setServerSymbol() {
	if strings.Count(g.Board, "O") == 1 {
		g.clientSymbol = SYMBOL_O
		g.serverSymbol = SYMBOL_X
	} else {
		g.clientSymbol = SYMBOL_X
		g.serverSymbol = SYMBOL_O
	}
}

func (g *Game) validateFirstInput() bool {
	countO := strings.Count(g.Board, "O")
	countX := strings.Count(g.Board, "X")
	return countO+countX <= 1
}

func (g *Game) validateBoard() bool {
	for _, c := range g.Board {
		if c != EMPTY && c != SYMBOL_O && c != SYMBOL_X {
			return false
		}
	}
	return true
}

func (g *Game) validateMove(next *Game) bool {
	moves := 0
	for i := 0; i < BOARD_LEN; i++ {
		if g.Board[i] == EMPTY && g.Board[i] != next.Board[i] {
			moves++
			if moves > 1 {
				return false
			}
			if next.Board[i] != g.clientSymbol {
				return false
			}
		} else if g.Board[i] != next.Board[i] {
			return false
		}
	}
	return true
}

func (g *Game) makeCounterMove() {
	emptyCells := g.findEmptyCells()
	randomIndex := emptyCells[g.randomGenerator.Intn(len(emptyCells))]
	g.Board = replaceAtIndex(g.Board, g.serverSymbol, randomIndex)
}

func (g *Game) findEmptyCells() []int {
	emptyCells := make([]int, 0)
	for i, cell := range g.Board {
		if cell == EMPTY {
			emptyCells = append(emptyCells, i)
		}
	}
	return emptyCells
}

func replaceAtIndex(board string, symbol byte, i int) string {
	out := []byte(board)
	out[i] = symbol
	return string(out)
}

func (g *Game) setStatus(idx int) {
	if g.Board[idx] == SYMBOL_O {
		g.Status = STATUS_O_WON
	} else {
		g.Status = STATUS_X_WON
	}
}

func (g *Game) checkRows() bool {
	for i := 0; i < 9; i += 3 {
		if g.Board[i] != EMPTY && g.Board[i] == g.Board[i+1] && g.Board[i+1] == g.Board[i+2] {
			g.setStatus(i)
			fmt.Printf("row %d\n", i)
			return true
		}
	}
	return false
}

func (g *Game) checkCols() bool {
	for i := 0; i < 3; i++ {
		if g.Board[i] != EMPTY && g.Board[i] == g.Board[i+3] && g.Board[i+3] == g.Board[i+6] {
			g.setStatus(i)
			fmt.Printf("col %d\n", i)
			return true
		}
	}
	return false
}

func (g *Game) checkDiagonal() bool {
	if g.Board[0] != EMPTY && g.Board[0] == g.Board[4] && g.Board[4] == g.Board[8] {
		g.setStatus(0)
		fmt.Printf("diag 1\n")
		return true
	}

	if g.Board[2] != EMPTY && g.Board[2] == g.Board[4] && g.Board[4] == g.Board[6] {
		g.setStatus(2)
		return true
	}
	return false
}

func (g *Game) checkDraw() bool {
	if emptycells := g.findEmptyCells(); len(emptycells) == 0 && g.Status == STATUS_RUNNING {
		g.Status = STATUS_DRAW
		return true
	}
	return false
}

func (g *Game) updateStatus() {
	_ = g.checkRows() ||
		g.checkCols() ||
		g.checkDiagonal() ||
		g.checkDraw()
}
