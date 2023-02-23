package main

import (
	"github.com/bengissimo/tictactoe/pkg/game"
)

func main() {
	gs := game.NewStore()

	gs.Router.Run("0.0.0.0:8080")
}
