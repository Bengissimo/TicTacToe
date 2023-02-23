package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func callCreateGame(router *gin.Engine, input string) (*Game, *httptest.ResponseRecorder, error) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/games", bytes.NewBufferString(input))

	router.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	game := &Game{}
	if err := json.Unmarshal(body, game); err != nil {
		return nil, nil, err
	}

	return game, w, nil
}

func callGetAllGames(router *gin.Engine) ([]Game, error) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/games", nil)

	router.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	games := []Game{}
	if err := json.Unmarshal(body, &games); err != nil {
		return nil, err
	}
	return games, nil
}

func callGetSingleGame(router *gin.Engine, id string) (*Game, *httptest.ResponseRecorder, error) {
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/games/%s", id)
	req, _ := http.NewRequest("GET", url, nil)

	router.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	game := &Game{}
	if err := json.Unmarshal(body, game); err != nil {
		return nil, w, err
	}

	return game, w, nil
}

func TestGame_CreateandGetGames(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantCreateCode    int
		wantLocation      string
		wantGetSingleCode int
		wantStatus        string
		wantLen           int
	}{
		{
			name:              "empty board",
			input:             `{"board":"---------"}`,
			wantCreateCode:    201,
			wantGetSingleCode: 200,
			wantStatus:        STATUS_RUNNING,
			wantLen:           1,
		},
		{
			name:              "board with 'O'",
			input:             `{"board":"---O-----"}`,
			wantCreateCode:    201,
			wantGetSingleCode: 200,
			wantStatus:        STATUS_RUNNING,
			wantLen:           2,
		},
		{
			name:              "board with 'X'",
			input:             `{"board":"----X----"}`,
			wantCreateCode:    201,
			wantGetSingleCode: 200,
			wantStatus:        STATUS_RUNNING,
			wantLen:           3,
		},
		{
			name:              "board with 'X'",
			input:             `{"board":"----X--"}`,
			wantCreateCode:    400,
			wantGetSingleCode: 404,
			wantStatus:        "",
			wantLen:           3,
		},
		{
			name:              "board with 'X'",
			input:             `{"board":"----XX--"}`,
			wantCreateCode:    400,
			wantGetSingleCode: 404,
			wantStatus:        "",
			wantLen:           3,
		},
		{
			name:              "board with 'a'",
			input:             `{"board":"----a----"}`,
			wantCreateCode:    400,
			wantGetSingleCode: 404,
			wantStatus:        "",
			wantLen:           3,
		},
	}
	store := NewStore()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//CreateGame
			game, w, err := callCreateGame(store.Router, tt.input)
			if err != nil {
				assert.Fail(t, "Game create failed")
			}

			assert.Equal(t, tt.wantCreateCode, w.Code)
			assert.Equal(t, tt.wantStatus, game.Status)
			if game.ID != uuid.Nil {
				tt.wantLocation = fmt.Sprintf("http://127.0.0.1:8080/api/v1/games/%s", game.ID.String())
				assert.Equal(t, tt.wantLocation, w.Header().Get("Location"))
			}

			//GetAllGames
			games, _ := callGetAllGames(store.Router)
			assert.Equal(t, tt.wantLen, len(games))

			//GetSingleGame
			game, w, err = callGetSingleGame(store.Router, game.ID.String())
			if err != nil {
				assert.Equal(t, tt.wantGetSingleCode, w.Code)
			}
			assert.Equal(t, tt.wantGetSingleCode, w.Code)
			assert.Equal(t, tt.wantStatus, game.Status)
		})
	}
}

func TestStore_DeleteGames(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantCode int
		wantBody string
	}{
		{
			name:     "valid",
			wantCode: 200,
			wantBody: `{"description":"Game successfully deleted"}`,
		},
		{
			name:     "invalid UUID",
			url:      "http://127.0.0.1:8080/api/v1/games/qweqwe",
			wantCode: 400,
			wantBody: `{"reason":"UUID cannot be parsed"}`,
		},
		{
			name:     "Wrong ID",
			url:      "http://127.0.0.1:8080/api/v1/games/00000000-0000-0000-0000-000000000000",
			wantCode: 404,
			wantBody: `{"reason":"Game not found"}`,
		},
	}

	store := NewStore()

	game, _, _ := callCreateGame(store.Router, `{"board":"---------"}`)
	tests[0].url = fmt.Sprintf("http://127.0.0.1:8080/api/v1/games/%s", game.ID.String())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", tt.url, nil)

			store.Router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestStore_MakeMove(t *testing.T) {
	tests := []struct {
		name       string
		board      string
		move       string
		wantCode   int
		wantStatus string
	}{
		{
			name:       "X wins",
			board:      "OXXOXOXO-",
			move:       `{"board":"OXXOXOXOX"}`,
			wantCode:   200,
			wantStatus: STATUS_X_WON,
		},
		{
			name:       "O wins",
			board:      "XOOXOXOO-",
			move:       `{"board":"XOOXOXOOX"}`,
			wantCode:   200,
			wantStatus: STATUS_O_WON,
		},
		{
			name:       "Draw",
			board:      "OXOXOXXO-",
			move:       `{"board":"OXOXOXXOX"}`,
			wantCode:   200,
			wantStatus: STATUS_DRAW,
		},
		{
			name:       "running",
			board:      "---X-O---",
			move:       `{"board":"---X-OX--"}`,
			wantCode:   200,
			wantStatus: STATUS_RUNNING,
		},
		{
			name:       "invalid input length",
			board:      "---X-O---",
			move:       `{"board":"----X-OX--"}`,
			wantCode:   400,
			wantStatus: "",
		},
		{
			name:       "invalid board input",
			board:      "---X-O---",
			move:       `{"board":"XXXX-O---"}`,
			wantCode:   400,
			wantStatus: "",
		},
	}

	store := NewStore()

	game, _, _ := callCreateGame(store.Router, `{"board":"---------"}`)
	url := fmt.Sprintf("http://127.0.0.1:8080/api/v1/games/%s", game.ID.String())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.Games[game.ID].Board = tt.board
			store.Games[game.ID].Status = STATUS_RUNNING

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", url, bytes.NewBufferString(tt.move))

			store.Router.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			gameOngoing := &Game{}
			if err := json.Unmarshal(body, gameOngoing); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantStatus, gameOngoing.Status)
		})
	}
}
