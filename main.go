package main

import (
	"github.com/bengissimo/tictactoe/pkg/game"
	"github.com/gin-gonic/gin"
)

func main() {
	gs := game.NewStore()

	router := gin.Default()

	router.GET("api/v1/games", gs.GetAllGames)
	router.GET("api/v1/games/:game_id", gs.GetSingleGame)
	router.POST("api/v1/games", gs.CreateGame)
	router.PUT("api/v1/games/:game_id", gs.MakeMove)
	router.DELETE("api/v1/games/:game_id", gs.DeleteGame)

	router.Run("localhost:8080")
}
