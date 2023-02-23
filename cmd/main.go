package main

import (
	"log"

	"github.com/bengissimo/tictactoe/pkg/game"
)

func main() {
	gs := game.NewStore()

	if err := gs.Router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal(err)
	}
}
