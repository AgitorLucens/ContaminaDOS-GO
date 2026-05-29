package querie

import (
	"be/internal/logic"
	"be/util/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func createAndStartGame(pool *logic.GamePool, suffix string, owner string) types.Game {
	g, _ := pool.CreateGame(*types.NewGame("Game"+suffix, owner, ""))

	pool.JoinGame("P2"+suffix, types.JoinRequest{GameId: g.GameId.String(), Player: "P2" + suffix, Password: ""})
	pool.JoinGame("P3"+suffix, types.JoinRequest{GameId: g.GameId.String(), Player: "P3" + suffix, Password: ""})
	pool.JoinGame("P4"+suffix, types.JoinRequest{GameId: g.GameId.String(), Player: "P4" + suffix, Password: ""})
	pool.JoinGame("P5"+suffix, types.JoinRequest{GameId: g.GameId.String(), Player: "P5" + suffix, Password: ""})

	pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(),
		Player: owner,
	})

	updated, _ := pool.GetGame(g.GameId.String(), types.GetGameRequest{})
	return *updated
}

func TestGetGame_Success(t *testing.T) {
	pool := logic.NewGamePool()
	g, _ := pool.CreateGame(*types.NewGame("TestGame", "Owner1", ""))

	r := setupRouter()
	r.GET("/api/games/:gameId", GetGame)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/games/"+g.GameId.String(), nil)
	req.Header.Set("Player", "Owner1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response types.Response
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Message != "Game Found" {
		t.Errorf("expected 'Game Found', got '%s'", response.Message)
	}
}

func TestGetGame_NotFound(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.GET("/api/games/:gameId", GetGame)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/games/"+uuid.New().String(), nil)
	req.Header.Set("Player", "Owner1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response types.Response
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Message != "Game Not Found" {
		t.Errorf("expected 'Game Not Found', got '%s'", response.Message)
	}
}

func TestGetRounds_Success(t *testing.T) {
	pool := logic.NewGamePool()
	g := createAndStartGame(pool, "Rounds", "Owner1")

	r := setupRouter()
	r.GET("/api/games/:gameId/rounds", GetRounds)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/games/"+g.GameId.String()+"/rounds", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetRounds_InvalidGameId(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.GET("/api/games/:gameId/rounds", GetRounds)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/games/invalid/rounds", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestShowRound_Success(t *testing.T) {
	pool := logic.NewGamePool()
	g := createAndStartGame(pool, "ShowRound", "Owner1")

	r := setupRouter()
	r.GET("/api/games/:gameId/rounds/:roundId", ShowRound)

	w := httptest.NewRecorder()
	url := "/api/games/" + g.GameId.String() + "/rounds/" + g.CurrentRound
	req, _ := http.NewRequest("GET", url, nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response types.Response
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Message != "Results found" {
		t.Errorf("expected 'Results found', got '%s'", response.Message)
	}
}

func TestShowRound_NotFound(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.GET("/api/games/:gameId/rounds/:roundId", ShowRound)

	w := httptest.NewRecorder()
	url := "/api/games/" + uuid.New().String() + "/rounds/" + uuid.New().String()
	req, _ := http.NewRequest("GET", url, nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestGameSearch_Success(t *testing.T) {
	pool := logic.NewGamePool()
	pool.CreateGame(*types.NewGame("Alpha", "Owner1", ""))

	r := setupRouter()
	r.GET("/api/games", GameSearch)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/games?name=Alpha&page=0&limit=50", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGameSearch_NoResults(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.GET("/api/games", GameSearch)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/games?page=0&limit=50", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetHello(t *testing.T) {
	r := setupRouter()
	r.GET("/hello", GetHello)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
