package main

import (
	"github.com/bengissimo/tictactoe/pkg/game"
)

func main() {
	gs := game.NewStore()

	gs.Router.Run("localhost:8080")
}
