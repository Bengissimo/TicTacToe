package game

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Store struct {
	Games map[uuid.UUID]*Game
}

func NewStore() *Store {
	return &Store{
		Games: make(map[uuid.UUID]*Game),
	}
}

func (s *Store) GetAllGames(c *gin.Context) {
	games := make([]*Game, 0)
	for _, g := range s.Games {
		games = append(games, g)
	}
	c.JSON(http.StatusOK, games)
}

func (s *Store) CreateGame(c *gin.Context) {
	newGame := Game{}

	if err := c.ShouldBindJSON(&newGame); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"reason": "Invalid board input"})
		return
	}

	newGame.ID = uuid.New()
	newGame.Status = STATUS_RUNNING
	//validateBoard
	s.Games[newGame.ID] = &newGame

	location := fmt.Sprintf("/api/v1/games/%s", newGame.ID.String())

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
		c.AbortWithStatusJSON(404, gin.H{"reason": "Invalid ID"})
		return nil
	}

	return game
}

func (s *Store) MakeMove(c *gin.Context) {
	game := s.getGameFromContext(c)
	if game == nil {
		return
	}

	newGame := Game{}
	if err := c.ShouldBindJSON(&newGame); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"reason": "Invalid board input"})
		return
	}
	game.Board = newGame.Board

	c.JSON(200, game)
}
