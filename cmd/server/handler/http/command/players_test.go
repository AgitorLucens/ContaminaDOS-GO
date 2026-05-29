package command

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"be/internal/logic"
	"be/util/types"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func createAndStartGame(pool *logic.GamePool, suffix string, owner string) (types.Game, types.Round) {
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
	rounds, _ := pool.GetRounds(g.GameId.String())
	return *updated, rounds[0]
}

func TestCreateGame_Success(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.POST("/api/games", CreateGame)

	body := types.CreateRequest{Name: "NewGame", Owner: "Owner1", Password: ""}
	w := httptest.NewRecorder()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/games", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCreateGame_MissingData(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.POST("/api/games", CreateGame)

	body := types.CreateRequest{Name: "", Owner: "", Password: ""}
	w := httptest.NewRecorder()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/games", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateGame_InvalidJSON(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.POST("/api/games", CreateGame)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/games", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateGame_Duplicate(t *testing.T) {
	pool := logic.NewGamePool()
	pool.CreateGame(*types.NewGame("ExistingGame", "Owner1", ""))

	r := setupRouter()
	r.POST("/api/games", CreateGame)

	body := types.CreateRequest{Name: "ExistingGame", Owner: "Owner2", Password: ""}
	w := httptest.NewRecorder()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/games", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestJoinGame_Success(t *testing.T) {
	pool := logic.NewGamePool()
	g, _ := pool.CreateGame(*types.NewGame("TestGame", "Owner1", ""))

	r := setupRouter()
	r.PUT("/api/games/:gameId", JoinGame)

	body := types.JoinRequest{GameId: g.GameId.String(), Player: "Player2", Password: ""}
	w := httptest.NewRecorder()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", "/api/games/"+g.GameId.String(), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestJoinGame_InvalidJSON(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.PUT("/api/games/:gameId", JoinGame)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/games/some-id", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGameStart_Success(t *testing.T) {
	pool := logic.NewGamePool()
	g, _ := pool.CreateGame(*types.NewGame("TestGame", "Owner1", ""))
	pool.JoinGame("P2", types.JoinRequest{GameId: g.GameId.String(), Player: "P2", Password: ""})
	pool.JoinGame("P3", types.JoinRequest{GameId: g.GameId.String(), Player: "P3", Password: ""})
	pool.JoinGame("P4", types.JoinRequest{GameId: g.GameId.String(), Player: "P4", Password: ""})
	pool.JoinGame("P5", types.JoinRequest{GameId: g.GameId.String(), Player: "P5", Password: ""})

	r := setupRouter()
	r.HEAD("/api/games/:gameId/start", GameStart)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("HEAD", "/api/games/"+g.GameId.String()+"/start", nil)
	req.Header.Set("Player", "Owner1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGameStart_Error(t *testing.T) {
	pool := logic.NewGamePool()
	g, _ := pool.CreateGame(*types.NewGame("TestGame", "Owner1", ""))

	r := setupRouter()
	r.HEAD("/api/games/:gameId/start", GameStart)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("HEAD", "/api/games/"+g.GameId.String()+"/start", nil)
	req.Header.Set("Player", "Owner1")
	r.ServeHTTP(w, req)

	if w.Header().Get("x-msg") == "" {
		t.Error("expected error header to be set")
	}
}

func TestProposeGroup_InvalidJSON(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.PATCH("/api/games/:gameId/rounds/:roundId", ProposeGroup)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/games/some/rounds/round", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestProposeGroup_GameNotFound(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.PATCH("/api/games/:gameId/rounds/:roundId", ProposeGroup)

	body := types.ProposeGroupRequest{Group: []string{"P1", "P2"}}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/games/"+uuid.New().String()+"/rounds/"+uuid.New().String(), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("player", "P1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestProposeGroup_Success(t *testing.T) {
	pool := logic.NewGamePool()
	g, round := createAndStartGame(pool, "Propose", "Owner1")

	r := setupRouter()
	r.PATCH("/api/games/:gameId/rounds/:roundId", ProposeGroup)

	leader := round.Leader
	group := []string{leader, "P2Propose"}
	body := types.ProposeGroupRequest{Group: group}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	url := "/api/games/" + g.GameId.String() + "/rounds/" + round.Id.String()
	req, _ := http.NewRequest("PATCH", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("player", leader)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestVoteGroup_InvalidJSON(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.POST("/api/games/:gameId/rounds/:roundId", VoteGroup)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/games/some/rounds/round", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestVoteGroup_Success(t *testing.T) {
	pool := logic.NewGamePool()
	g, round := createAndStartGame(pool, "Vote", "Owner1")

	leader := round.Leader
	group := []string{leader, "P2Vote"}
	pool.ProposeGroup(g.GameId.String(), round.Id.String(), leader, group)

	r := setupRouter()
	r.POST("/api/games/:gameId/rounds/:roundId", VoteGroup)

	body := types.VoteRequest{Vote: true}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	url := "/api/games/" + g.GameId.String() + "/rounds/" + round.Id.String()
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("player", "Owner1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestSubmitAction_InvalidJSON(t *testing.T) {
	logic.NewGamePool()
	r := setupRouter()
	r.PUT("/api/games/:gameId/rounds/:roundId", SubmitAction)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/games/some/rounds/round", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSubmitAction_Success(t *testing.T) {
	pool := logic.NewGamePool()
	g, round := createAndStartGame(pool, "Action", "Owner1")

	leader := round.Leader
	group := []string{leader, "P2Action"}
	pool.ProposeGroup(g.GameId.String(), round.Id.String(), leader, group)

	players := []string{"Owner1", "P2Action", "P3Action", "P4Action", "P5Action"}
	for _, p := range players {
		pool.VoteGroup(g.GameId.String(), round.Id.String(), p, true)
	}

	r := setupRouter()
	r.PUT("/api/games/:gameId/rounds/:roundId", SubmitAction)

	body := types.ActionRequest{Action: true}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	url := "/api/games/" + g.GameId.String() + "/rounds/" + round.Id.String()
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("player", leader)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
