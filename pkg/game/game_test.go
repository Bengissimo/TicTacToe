package game

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame_setServerSymbol(t *testing.T) {
	tests := []struct {
		name        string
		board       string
		expectedSymbol byte
	}{
		{
			name:        "empty board",
			board:       "---------",
			expectedSymbol: SYMBOL_O,
		},
		{
			name:        "board with 'O'",
			board:       "---O-----",
			expectedSymbol: SYMBOL_X,
		},
		{
			name:        "board with 'X'",
			board:       "---X-----",
			expectedSymbol: SYMBOL_O,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board: tt.board,
			}

			g.setServerSymbol()

			/*if g.serverSymbol != tt.expectedSym {
				t.Errorf("Expected server symbol to be %c, but got %c", tt.expectedSym, g.serverSymbol)
			}*/
			assert.Equal(t, g.serverSymbol, tt.expectedSymbol)
		})
	}
}

func TestGame_findEmptyCells(t *testing.T) {
	tests := []struct {
		name     string
		board    string
		expected []int
	}{
		{
			name:     "empty board",
			board:    "---------",
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:     "single empty place",
			board:    "XOXOXOX-O",
			expected: []int{7},
		},
		{
			name:     "full",
			board:    "XOXOXOXOO",
			expected: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board: tt.board,
			}
			assert.Equal(t, g.findEmptyCells(), tt.expected)
		})
	}
}

func TestGame_makeCounterMove(t *testing.T) {
	tests := []struct {
		name          string
		board         string
		expectedBoard string
	}{
		{
			name:          "empty board",
			board:         "---------",
			expectedBoard: "O--------",
		},
		{
			name:          "first is occupied",
			board:         "X--------",
			expectedBoard: "X--O-----",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board:           tt.board,
				randomGenerator: rand.New(rand.NewSource(0)),
				serverSymbol:    SYMBOL_O,
			}
			g.makeCounterMove()
			assert.Equal(t, g.Board, tt.expectedBoard)
		})
	}
}

func TestGame_validateFirstInput(t *testing.T) {
	tests := []struct {
		name     string
		board    string
		expected bool
	}{
		{
			name:     "empty input",
			board:    "---------",
			expected: true,
		},
		{
			name:     "valid X",
			board:    "X--------",
			expected: true,
		},
		{
			name:     "valid O",
			board:    "O--------",
			expected: true,
		},
		{
			name:     "invalid O",
			board:    "O--O-----",
			expected: false,
		},
		{
			name:     "invalid X",
			board:    "X--X-----",
			expected: false,
		},
		{
			name:     "invalid OX",
			board:    "O--X-----",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board: tt.board,
			}
			assert.Equal(t, g.validateFirstInput(), tt.expected)
		})
	}
}

func TestGame_validateBoard(t *testing.T) {
	tests := []struct {
		name     string
		board    string
		expected bool
	}{
		{
			name:     "empty input",
			board:    "---------",
			expected: true,
		},
		{
			name:     "valid X",
			board:    "X--------",
			expected: true,
		},
		{
			name:     "valid O",
			board:    "O--------",
			expected: true,
		},
		{
			name:     "valid XX",
			board:    "--X-X----",
			expected: true,
		},
		{
			name:     "valid OO",
			board:    "O-O------",
			expected: true,
		},
		{
			name:     "valid OX",
			board:    "O-X-O----",
			expected: true,
		},
		{
			name:     "invalid 1",
			board:    "Oa-------",
			expected: false,
		},
		{
			name:     "invalid 2",
			board:    "-a-------",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board: tt.board,
			}
			assert.Equal(t, g.validateBoard(), tt.expected)
		})
	}
}

func TestGame_validateMove(t *testing.T) {
	tests := []struct {
		name     string
		board    string
		next     string
		expected bool
	}{
		{
			name:     "valid ",
			board:    "--O------",
			next:     "--O---X--",
			expected: true,
		},
		{
			name:     "invalid same place",
			board:    "--X---O--",
			next:     "--X---X--",
			expected: false,
		},
		{
			name:     "valid full",
			board:    "OXOXOXOX-",
			next:     "OXOXOXOXX",
			expected: true,
		},
		{
			name:     "invalid full",
			board:    "OXOXOXOX-",
			next:     "OXOXOXOXX",
			expected: true,
		},
		{
			name:     "invalid move >1",
			board:    "--OX--O--",
			next:     "--OXXXO--",
			expected: false,
		},
		{
			name:     "invalid serversymbol",
			board:    "--O------",
			next:     "--O---O--",
			expected: false,
		},
		{
			name:     "invalid symbol",
			board:    "--O------",
			next:     "--O---a--",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board:        tt.board,
				clientSymbol: SYMBOL_X,
			}
			next := &Game{
				Board: tt.next,
			}
			assert.Equal(t, g.validateMove(next), tt.expected)
		})
	}
}

func TestGame_checkRows(t *testing.T) {
	tests := []struct {
		name     string
		board    string
		expected bool
	}{
		{
			name:     "empty input",
			board:    "---------",
			expected: false,
		},
		{
			name:     "valid 1",
			board:    "OOO------",
			expected: true,
		},
		{
			name:     "valid 2",
			board:    "---XXX---",
			expected: true,
		},
		{
			name:     "valid 3",
			board:    "------OOO",
			expected: true,
		},
		{
			name:     "valid 4",
			board:    "-XX---OOO",
			expected: true,
		},
		{
			name:     "invalid 1",
			board:    "------OXO",
			expected: false,
		},
		{
			name:     "invalid 2",
			board:    "--XX--O-O",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board: tt.board,
			}
			assert.Equal(t, g.checkRows(), tt.expected)
		})
	}
}

func TestGame_checkCols(t *testing.T) {
	tests := []struct {
		name     string
		board    string
		expected bool
	}{
		{
			name:     "empty input",
			board:    "---------",
			expected: false,
		},
		{
			name:     "valid 1",
			board:    "O--O--O--",
			expected: true,
		},
		{
			name:     "valid 2",
			board:    "-X--X--X-",
			expected: true,
		},
		{
			name:     "valid 3",
			board:    "--O--O--O",
			expected: true,
		},
		{
			name:     "valid 4",
			board:    "-X--XOOX-",
			expected: true,
		},
		{
			name:     "invalid 1",
			board:    "O--X--O--",
			expected: false,
		},
		{
			name:     "invalid 2",
			board:    "--O-OX--O",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board: tt.board,
			}
			assert.Equal(t, g.checkCols(), tt.expected)
		})
	}
}

func TestGame_checkDiagonal(t *testing.T) {
	tests := []struct {
		name     string
		board    string
		expected bool
	}{
		{
			name:     "empty input",
			board:    "---------",
			expected: false,
		},
		{
			name:     "valid 1",
			board:    "O---O---O",
			expected: true,
		},
		{
			name:     "valid 2",
			board:    "--X-X-X--",
			expected: true,
		},
		{
			name:     "valid 3",
			board:    "--XOXOX--",
			expected: true,
		},
		{
			name:     "invalid 1",
			board:    "O---O----",
			expected: false,
		},
		{
			name:     "invalid 2",
			board:    "--XXOOX--",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board: tt.board,
			}
			assert.Equal(t, g.checkDiagonal(), tt.expected)
		})
	}
}

func TestGame_checkDraw(t *testing.T) {
	tests := []struct {
		name     string
		board    string
		status   string
		expected bool
	}{
		{
			name:     "empty input",
			board:    "---------",
			status:   STATUS_RUNNING,
			expected: false,
		},
		{
			name:     "valid",
			board:    "OXOXOXOXO",
			status:   STATUS_RUNNING,
			expected: true,
		},
		{
			name:     "invalid",
			board:    "OXOXOXOXO",
			status:   STATUS_O_WON,
			expected: false,
		},
		{
			name:     "invalid 2",
			board:    "XOXOOOXOO",
			status:   STATUS_O_WON,
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board:  tt.board,
				Status: tt.status,
			}
			assert.Equal(t, g.checkDraw(), tt.expected)
		})
	}
}
