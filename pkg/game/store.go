package game

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Store struct {
	Games           map[uuid.UUID]*Game
	randomGenerator *rand.Rand
	Router          *gin.Engine
}

func NewStore() *Store {
	gs := &Store{
		Games:           make(map[uuid.UUID]*Game),
		randomGenerator: rand.New(rand.NewSource(time.Now().UnixNano())),
		Router:          gin.Default(),
	}

	gs.Router.GET("api/v1/games", gs.GetAllGames)
	gs.Router.GET("api/v1/games/:game_id", gs.GetSingleGame)
	gs.Router.POST("api/v1/games", gs.CreateGame)
	gs.Router.PUT("api/v1/games/:game_id", gs.MakeMove)
	gs.Router.DELETE("api/v1/games/:game_id", gs.DeleteGame)

	return gs
}

func (s *Store) GetAllGames(c *gin.Context) {
	games := make([]*Game, 0)

	for _, g := range s.Games {
		games = append(games, g)
	}

	c.JSON(http.StatusOK, games)
}

func (s *Store) CreateGame(c *gin.Context) {
	newGame := Game{
		randomGenerator: s.randomGenerator,
	}

	if err := c.ShouldBindJSON(&newGame); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"reason": "Invalid input length"})
		return
	}

	newGame.Board = strings.ToUpper(newGame.Board)

	if !newGame.validateFirstInput() || !newGame.validateBoard() {
		c.AbortWithStatusJSON(400, gin.H{"reason": "Invalid board input"})
		return
	}

	newGame.ID = uuid.New()
	newGame.Status = STATUS_RUNNING

	s.Games[newGame.ID] = &newGame

	newGame.setServerSymbol()
	newGame.makeCounterMove()

	location := fmt.Sprintf("http://127.0.0.1:8080/api/v1/games/%s", newGame.ID.String())
	c.Header("Location", location)

	c.JSON(201, newGame)
}

func (s *Store) GetSingleGame(c *gin.Context) {
	game := s.getGameFromContext(c)
	if game == nil {
		return
	}

	c.JSON(200, game)
}

func (s *Store) DeleteGame(c *gin.Context) {
	game := s.getGameFromContext(c)
	if game == nil {
		return
	}

	delete(s.Games, game.ID)

	c.JSON(200, gin.H{"description": "Game successfully deleted"})
}

func (s *Store) getGameFromContext(c *gin.Context) *Game {
	id := c.Param("game_id")
	gameID, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"reason": "UUID cannot be parsed"})
		return nil
	}

	game, ok := s.Games[gameID]
	if !ok {
		c.AbortWithStatusJSON(404, gin.H{"reason": "Game not found"})
		return nil
	}

	return game
}

func (s *Store) MakeMove(c *gin.Context) {
	game := s.getGameFromContext(c)
	if game == nil {
		return
	}

	newGame := &Game{}

	if err := c.ShouldBindJSON(newGame); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"reason": "Invalid input length"})
		return
	}

	newGame.Board = strings.ToUpper(newGame.Board)

	if !newGame.validateBoard() || !game.validateMove((newGame)) {
		c.AbortWithStatusJSON(400, gin.H{"reason": "Invalid board input"})
		return
	}

	game.Board = newGame.Board
	fmt.Printf("%s\n", game.Board)

	game.updateStatus()
	if game.Status != STATUS_RUNNING {
		c.JSON(200, game)
		return
	}

	game.makeCounterMove()

	game.updateStatus()
	if game.Status != STATUS_RUNNING {
		c.JSON(200, game)
		return
	}

	c.JSON(200, game)
}
