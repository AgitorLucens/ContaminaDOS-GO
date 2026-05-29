package logic

import (
	"be/util/types"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func resetPool() {
	gamePoolInstance = nil
	once = syncOnceReset()
}

func syncOnceReset() sync.Once {
	return sync.Once{}
}

func setupGamePool() *GamePool {
	resetPool()
	return NewGamePool()
}

func createTestGame(pool *GamePool, name, owner string) types.Game {
	g, _ := pool.CreateGame(*types.NewGame(name, owner, ""))
	return g
}

func addTestPlayers(pool *GamePool, g *types.Game, players []string) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	id := g.GameId
	game := pool.gamesUUID[id]
	for _, p := range players {
		game.Players = append(game.Players, p)
	}
	pool.gamesUUID[id] = game
	pool.games[game.Name] = game
	*g = game
}

func TestCreateGame_Success(t *testing.T) {
	pool := setupGamePool()
	g, err := pool.CreateGame(*types.NewGame("TestGame", "Owner1", ""))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if g.Name != "TestGame" {
		t.Errorf("expected name TestGame, got %s", g.Name)
	}
	if g.Owner != "Owner1" {
		t.Errorf("expected owner Owner1, got %s", g.Owner)
	}
	if g.Status != types.StatusLobby {
		t.Errorf("expected status lobby, got %s", g.Status)
	}
	if len(g.Players) != 1 || g.Players[0] != "Owner1" {
		t.Errorf("expected owner in players, got %v", g.Players)
	}
	if g.GameId == uuid.Nil {
		t.Error("expected non-nil game ID")
	}
}

func TestCreateGame_Duplicate(t *testing.T) {
	pool := setupGamePool()
	pool.CreateGame(*types.NewGame("TestGame", "Owner1", ""))
	_, err := pool.CreateGame(*types.NewGame("TestGame", "Owner2", ""))
	if err == nil {
		t.Fatal("expected error for duplicate game, got nil")
	}
	if err.Error() != "Game already exists" {
		t.Errorf("expected 'Game already exists', got '%s'", err.Error())
	}
}

func TestCreateGame_WithPassword(t *testing.T) {
	pool := setupGamePool()
	g, _ := pool.CreateGame(*types.NewGame("TestGame", "Owner1", "secret123"))
	if g.PWD != "secret123" {
		t.Errorf("expected PWD secret123, got %s", g.PWD)
	}
}

func TestSearchGamesByNameAndStatus_Empty(t *testing.T) {
	pool := setupGamePool()
	results := pool.SearchGamesByNameAndStatus("", "", 0, 50)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchGamesByNameAndStatus_Found(t *testing.T) {
	pool := setupGamePool()
	pool.CreateGame(*types.NewGame("Alpha", "Owner1", ""))
	pool.CreateGame(*types.NewGame("Beta", "Owner2", ""))
	results := pool.SearchGamesByNameAndStatus("Alpha", "", 0, 50)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "Alpha" {
		t.Errorf("expected Alpha, got %s", results[0].Name)
	}
}

func TestSearchGamesByNameAndStatus_ByStatus(t *testing.T) {
	pool := setupGamePool()
	pool.CreateGame(*types.NewGame("Alpha", "Owner1", ""))
	results := pool.SearchGamesByNameAndStatus("", "lobby", 0, 50)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchGamesByNameAndStatus_Pagination(t *testing.T) {
	pool := setupGamePool()
	pool.CreateGame(*types.NewGame("Game1", "O1", ""))
	pool.CreateGame(*types.NewGame("Game2", "O2", ""))
	pool.CreateGame(*types.NewGame("Game3", "O3", ""))
	results := pool.SearchGamesByNameAndStatus("", "", 1, 1)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchGamesByNameAndStatus_PageOutOfBounds(t *testing.T) {
	pool := setupGamePool()
	pool.CreateGame(*types.NewGame("Game1", "O1", ""))
	results := pool.SearchGamesByNameAndStatus("", "", 5, 10)
	if len(results) != 0 {
		t.Errorf("expected 0 results for out-of-bounds page, got %d", len(results))
	}
}

func TestSearchGamesByNameAndStatus_NoMatch(t *testing.T) {
	pool := setupGamePool()
	pool.CreateGame(*types.NewGame("Game1", "O1", ""))
	results := pool.SearchGamesByNameAndStatus("NonExistent", "", 0, 50)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestGetGame_Success(t *testing.T) {
	pool := setupGamePool()
	g, _ := pool.CreateGame(*types.NewGame("TestGame", "Owner1", ""))
	found, err := pool.GetGame(g.GameId.String(), types.GetGameRequest{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Name != "TestGame" {
		t.Errorf("expected name TestGame, got %s", found.Name)
	}
}

func TestGetGame_InvalidID(t *testing.T) {
	pool := setupGamePool()
	_, err := pool.GetGame("invalid-uuid", types.GetGameRequest{})
	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}
}

func TestGetGame_NotFound(t *testing.T) {
	pool := setupGamePool()
	_, err := pool.GetGame(uuid.New().String(), types.GetGameRequest{})
	if err == nil {
		t.Fatal("expected error for not found game, got nil")
	}
	if err.Error() != "Game not found" {
		t.Errorf("expected 'Game not found', got '%s'", err.Error())
	}
}

func TestJoinGame_Success(t *testing.T) {
	pool := setupGamePool()
	g, _ := pool.CreateGame(*types.NewGame("TestGame", "Owner1", ""))
	req := types.JoinRequest{
		GameId:   g.GameId.String(),
		Player:   "Player2",
		Password: "",
	}
	joined, err := pool.JoinGame("Player2", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	found := false
	for _, p := range joined.Players {
		if p == "Player2" {
			found = true
		}
	}
	if !found {
		t.Error("expected Player2 in game players")
	}
}

func TestJoinGame_InvalidGameID(t *testing.T) {
	pool := setupGamePool()
	req := types.JoinRequest{GameId: "bad", Player: "P1", Password: ""}
	_, err := pool.JoinGame("P1", req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestJoinGame_GameNotFound(t *testing.T) {
	pool := setupGamePool()
	req := types.JoinRequest{GameId: uuid.New().String(), Player: "P1", Password: ""}
	_, err := pool.JoinGame("P1", req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "Game not found" {
		t.Errorf("expected 'Game not found', got '%s'", err.Error())
	}
}

func TestJoinGame_WrongPassword(t *testing.T) {
	pool := setupGamePool()
	g, _ := pool.CreateGame(*types.NewGame("TestGame", "Owner1", "secret"))
	req := types.JoinRequest{GameId: g.GameId.String(), Player: "P1", Password: "wrong"}
	_, err := pool.JoinGame("P1", req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "Invalid game password" {
		t.Errorf("expected 'Invalid game password', got '%s'", err.Error())
	}
}

func TestJoinGame_DuplicatePlayer(t *testing.T) {
	pool := setupGamePool()
	g, _ := pool.CreateGame(*types.NewGame("TestGame", "Owner1", ""))
	req := types.JoinRequest{GameId: g.GameId.String(), Player: "Owner1", Password: ""}
	_, err := pool.JoinGame("Owner1", req)
	if err == nil {
		t.Fatal("expected error for duplicate player, got nil")
	}
}

func TestStartGame_Success(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	players := []string{"P2", "P3", "P4", "P5"}
	addTestPlayers(pool, &g, players)

	req := types.StartGameRequest{
		GameId:   g.GameId.String(),
		Player:   "Owner1",
		Password: "",
	}
	started, err := pool.StartGame(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if started.Status != types.StatusRounds {
		t.Errorf("expected status rounds, got %s", started.Status)
	}
	if started.CurrentRound == "0000000000000000000000000" {
		t.Error("expected round ID to be set after start")
	}
	if len(started.Enemies) == 0 {
		t.Error("expected enemies to be assigned")
	}
}

func TestStartGame_InvalidID(t *testing.T) {
	pool := setupGamePool()
	req := types.StartGameRequest{GameId: "bad", Player: "O1", Password: ""}
	_, err := pool.StartGame(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestStartGame_GameNotFound(t *testing.T) {
	pool := setupGamePool()
	req := types.StartGameRequest{GameId: uuid.New().String(), Player: "O1", Password: ""}
	_, err := pool.StartGame(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestStartGame_NotEnoughPlayers(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3"})

	req := types.StartGameRequest{GameId: g.GameId.String(), Player: "Owner1", Password: ""}
	_, err := pool.StartGame(req)
	if err == nil {
		t.Fatal("expected error for not enough players, got nil")
	}
	if err.Error() != "Need 5 players to start" {
		t.Errorf("expected 'Need 5 players to start', got '%s'", err.Error())
	}
}

func TestStartGame_AlreadyStarted(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	pool.mutex.Lock()
	game := pool.gamesUUID[g.GameId]
	game.SetStatus(types.StatusRounds)
	pool.gamesUUID[g.GameId] = game
	pool.games[game.Name] = game
	pool.mutex.Unlock()

	req := types.StartGameRequest{GameId: g.GameId.String(), Player: "Owner1", Password: ""}
	_, err := pool.StartGame(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "Game already started" {
		t.Errorf("expected 'Game already started', got '%s'", err.Error())
	}
}

func TestStartGame_NotOwner(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	req := types.StartGameRequest{GameId: g.GameId.String(), Player: "P2", Password: ""}
	_, err := pool.StartGame(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "You are not the owner of this game" {
		t.Errorf("expected 'You are not the owner of this game', got '%s'", err.Error())
	}
}

func TestStartGame_WrongPassword(t *testing.T) {
	pool := setupGamePool()
	g, _ := pool.CreateGame(*types.NewGame("TestGame", "Owner1", "secret"))
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	req := types.StartGameRequest{GameId: g.GameId.String(), Player: "Owner1", Password: "wrong"}
	_, err := pool.StartGame(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetRounds_Success(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	req := types.StartGameRequest{GameId: g.GameId.String(), Player: "Owner1", Password: ""}
	pool.StartGame(req)

	rounds, err := pool.GetRounds(g.GameId.String())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(rounds) != 1 {
		t.Errorf("expected 1 round, got %d", len(rounds))
	}
}

func TestGetRounds_InvalidGameID(t *testing.T) {
	pool := setupGamePool()
	_, err := pool.GetRounds("invalid-uuid")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetRounds_NoRounds(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	rounds, _ := pool.GetRounds(g.GameId.String())
	if len(rounds) != 0 {
		t.Errorf("expected 0 rounds, got %d", len(rounds))
	}
}

func TestShowRound_Success(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	rounds, _ := pool.GetRounds(g.GameId.String())
	if len(rounds) == 0 {
		t.Fatal("expected at least 1 round")
	}

	round, err := pool.ShowRound(g.GameId.String(), started.CurrentRound)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if round.Id.String() != started.CurrentRound {
		t.Errorf("expected round ID %s, got %s", started.CurrentRound, round.Id.String())
	}
}

func TestShowRound_InvalidGameID(t *testing.T) {
	pool := setupGamePool()
	_, err := pool.ShowRound("invalid", uuid.New().String())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestShowRound_InvalidRoundID(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	_, err := pool.ShowRound(g.GameId.String(), "invalid")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestShowRound_NotFound(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	_, err := pool.ShowRound(g.GameId.String(), uuid.New().String())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "Round not found" {
		t.Errorf("expected 'Round not found', got '%s'", err.Error())
	}
}

func TestProposeGroup_Success(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	rounds, _ := pool.GetRounds(g.GameId.String())
	leader := rounds[0].Leader

	round, err := pool.ProposeGroup(g.GameId.String(), started.CurrentRound, leader, []string{leader, "P2"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if round.Status != types.RoundStatusVoting {
		t.Errorf("expected status voting, got %s", round.Status)
	}
	if len(round.Group) != 2 {
		t.Errorf("expected 2 group members, got %d", len(round.Group))
	}
}

func TestProposeGroup_InvalidGameID(t *testing.T) {
	pool := setupGamePool()
	_, err := pool.ProposeGroup("invalid", uuid.New().String(), "P1", []string{"P1"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestProposeGroup_InvalidRoundID(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	_, err := pool.ProposeGroup(g.GameId.String(), "invalid", "P1", []string{"P1"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestProposeGroup_GameNotFound(t *testing.T) {
	pool := setupGamePool()
	_, err := pool.ProposeGroup(uuid.New().String(), uuid.New().String(), "P1", []string{"P1"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "Game not found" {
		t.Errorf("expected 'Game not found', got '%s'", err.Error())
	}
}

func TestProposeGroup_NotLeader(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	_, err := pool.ProposeGroup(g.GameId.String(), started.CurrentRound, "NotLeader", []string{"P2", "P3"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "You are not the Leader" {
		t.Errorf("expected 'You are not the Leader', got '%s'", err.Error())
	}
}

func TestProposeGroup_PlayerNotInGame(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	rounds, _ := pool.GetRounds(g.GameId.String())
	leader := rounds[0].Leader

	_, err := pool.ProposeGroup(g.GameId.String(), started.CurrentRound, leader, []string{"NonExistent", "P2"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "Group member is not part of the game" {
		t.Errorf("expected 'Group member is not part of the game', got '%s'", err.Error())
	}
}

func TestProposeGroup_WrongGroupSize(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	rounds, _ := pool.GetRounds(g.GameId.String())
	leader := rounds[0].Leader

	_, err := pool.ProposeGroup(g.GameId.String(), started.CurrentRound, leader, []string{leader, "P2", "P3"})
	if err == nil {
		t.Fatal("expected error for wrong group size, got nil")
	}
}

func TestProposeGroup_NotWaitingOnLeader(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	pool.mutex.Lock()
	for id, r := range pool.roundsUUID {
		if r.Id.String() == started.CurrentRound {
			r.Status = types.RoundStatusVoting
			pool.roundsUUID[id] = r
		}
	}
	pool.mutex.Unlock()

	rounds, _ := pool.GetRounds(g.GameId.String())
	leader := rounds[0].Leader

	_, err := pool.ProposeGroup(g.GameId.String(), started.CurrentRound, leader, []string{leader, "P2"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "It is not the time for proposing groups" {
		t.Errorf("expected 'It is not the time for proposing groups', got '%s'", err.Error())
	}
}

func TestVoteGroup_Success(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	rounds, _ := pool.GetRounds(g.GameId.String())
	leader := rounds[0].Leader

	pool.ProposeGroup(g.GameId.String(), started.CurrentRound, leader, []string{leader, "P2"})

	round, err := pool.VoteGroup(g.GameId.String(), started.CurrentRound, "Owner1", true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(round.Votes) != 1 {
		t.Errorf("expected 1 vote, got %d", len(round.Votes))
	}
}

func TestVoteGroup_InvalidGameID(t *testing.T) {
	pool := setupGamePool()
	_, err := pool.VoteGroup("invalid", uuid.New().String(), "P1", true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestVoteGroup_PlayerNotInGame(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	_, err := pool.VoteGroup(g.GameId.String(), started.CurrentRound, "NonExistent", true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "Player is not part of the game" {
		t.Errorf("expected 'Player is not part of the game', got '%s'", err.Error())
	}
}

func TestVoteGroup_AlreadyVoted(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	rounds, _ := pool.GetRounds(g.GameId.String())
	leader := rounds[0].Leader

	pool.ProposeGroup(g.GameId.String(), started.CurrentRound, leader, []string{leader, "P2"})
	pool.VoteGroup(g.GameId.String(), started.CurrentRound, "Owner1", true)

	_, err := pool.VoteGroup(g.GameId.String(), started.CurrentRound, "Owner1", true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "You have already voted" {
		t.Errorf("expected 'You have already voted', got '%s'", err.Error())
	}
}

func TestVoteGroup_NotVotingStatus(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	_, err := pool.VoteGroup(g.GameId.String(), started.CurrentRound, "Owner1", true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "It is not the time for voting" {
		t.Errorf("expected 'It is not the time for voting', got '%s'", err.Error())
	}
}

func TestVoteGroup_VoteRejectedAdvancesPhase(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	rounds, _ := pool.GetRounds(g.GameId.String())
	leader := rounds[0].Leader

	pool.ProposeGroup(g.GameId.String(), started.CurrentRound, leader, []string{leader, "P2"})

	players := []string{"Owner1", "P2", "P3", "P4", "P5"}
	for i, p := range players {
		vote := true
		if i > 1 {
			vote = false
		}
		pool.VoteGroup(g.GameId.String(), started.CurrentRound, p, vote)
	}

	round, _ := pool.ShowRound(g.GameId.String(), started.CurrentRound)
	if round.Phase != types.RoundPhaseVote2 {
		t.Errorf("expected phase vote2, got %s", round.Phase)
	}
	if round.Status != types.RoundStatusWaitingOnLeader {
		t.Errorf("expected status waiting-on-leader, got %s", round.Status)
	}
}

func TestSubmitAction_Success(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	rounds, _ := pool.GetRounds(g.GameId.String())
	leader := rounds[0].Leader
	group := []string{leader, "P2"}

	pool.ProposeGroup(g.GameId.String(), started.CurrentRound, leader, group)

	players := []string{"Owner1", "P2", "P3", "P4", "P5"}
	for _, p := range players {
		pool.VoteGroup(g.GameId.String(), started.CurrentRound, p, true)
	}

	round, err := pool.SubmitAction(g.GameId.String(), started.CurrentRound, leader, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(round.Actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(round.Actions))
	}
}

func TestSubmitAction_InvalidGameID(t *testing.T) {
	pool := setupGamePool()
	_, err := pool.SubmitAction("invalid", uuid.New().String(), "P1", true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSubmitAction_PlayerNotInGame(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	_, err := pool.SubmitAction(g.GameId.String(), started.CurrentRound, "NonExistent", true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSubmitAction_PlayerNotInGroup(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	updated, _ := pool.GetGame(g.GameId.String(), types.GetGameRequest{})
	rounds, _ := pool.GetRounds(g.GameId.String())
	leader := rounds[0].Leader
	pool.ProposeGroup(g.GameId.String(), updated.CurrentRound, leader, []string{leader, "P2"})

	players := []string{"Owner1", "P2", "P3", "P4", "P5"}
	for _, p := range players {
		pool.VoteGroup(g.GameId.String(), updated.CurrentRound, p, true)
	}

	_, err := pool.SubmitAction(g.GameId.String(), updated.CurrentRound, "P3", true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "You cannot contribute in this round" {
		t.Errorf("expected 'You cannot contribute in this round', got '%s'", err.Error())
	}
}

func TestSubmitAction_NotWaitingOnGroup(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	updated, _ := pool.GetGame(g.GameId.String(), types.GetGameRequest{})

	_, err := pool.SubmitAction(g.GameId.String(), updated.CurrentRound, "P3", true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "You cannot contribute in this round" {
		t.Errorf("expected 'You cannot contribute in this round', got '%s'", err.Error())
	}
}

func TestSubmitAction_AllActionsSubmitted_RoundEnds(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	rounds, _ := pool.GetRounds(g.GameId.String())
	leader := rounds[0].Leader
	group := []string{leader, "P2"}

	pool.ProposeGroup(g.GameId.String(), started.CurrentRound, leader, group)

	players := []string{"Owner1", "P2", "P3", "P4", "P5"}
	for _, p := range players {
		pool.VoteGroup(g.GameId.String(), started.CurrentRound, p, true)
	}

	pool.SubmitAction(g.GameId.String(), started.CurrentRound, leader, true)
	round, _ := pool.SubmitAction(g.GameId.String(), started.CurrentRound, "P2", true)

	if round.Status != types.RoundStatusEnded {
		t.Errorf("expected status ended, got %s", round.Status)
	}
	if round.Result != types.RoundResultCitizens {
		t.Errorf("expected result citizens, got %s", round.Result)
	}
}

func TestSubmitAction_SabotageRound(t *testing.T) {
	pool := setupGamePool()
	g := createTestGame(pool, "TestGame", "Owner1")
	addTestPlayers(pool, &g, []string{"P2", "P3", "P4", "P5"})

	started, _ := pool.StartGame(types.StartGameRequest{
		GameId: g.GameId.String(), Player: "Owner1", Password: "",
	})

	rounds, _ := pool.GetRounds(g.GameId.String())
	leader := rounds[0].Leader
	group := []string{leader, "P2"}

	pool.ProposeGroup(g.GameId.String(), started.CurrentRound, leader, group)

	players := []string{"Owner1", "P2", "P3", "P4", "P5"}
	for _, p := range players {
		pool.VoteGroup(g.GameId.String(), started.CurrentRound, p, true)
	}

	pool.SubmitAction(g.GameId.String(), started.CurrentRound, leader, true)
	round, _ := pool.SubmitAction(g.GameId.String(), started.CurrentRound, "P2", false)

	if round.Result != types.RoundResultEnemies {
		t.Errorf("expected result enemies, got %s", round.Result)
	}
}

func TestAddPlayer_Success(t *testing.T) {
	g := types.NewGame("Test", "Owner1", "")
	g.Players = append(g.Players, "Owner1")
	err := AddPlayer("P2", g)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(g.Players) != 2 {
		t.Errorf("expected 2 players, got %d", len(g.Players))
	}
}

func TestAddPlayer_Duplicate(t *testing.T) {
	g := types.NewGame("Test", "Owner1", "")
	g.Players = append(g.Players, "Owner1")
	err := AddPlayer("Owner1", g)
	if err == nil {
		t.Fatal("expected error for duplicate, got nil")
	}
}

func TestGetEnemiesCountAtStart(t *testing.T) {
	tests := []struct {
		players  int
		expected int
	}{
		{5, 2},
		{6, 2},
		{7, 3},
		{8, 3},
		{9, 3},
		{10, 4},
	}
	for _, tt := range tests {
		result := getEnemiesCountAtStart(tt.players)
		if result != tt.expected {
			t.Errorf("getEnemiesCountAtStart(%d) = %d, expected %d", tt.players, result, tt.expected)
		}
	}
}

func TestGetExpectedGroupSize(t *testing.T) {
	tests := []struct {
		players    int
		roundCount int
		expected   int
	}{
		{5, 1, 2},
		{5, 2, 3},
		{6, 1, 2},
		{6, 3, 4},
		{7, 1, 2},
		{8, 1, 3},
		{10, 5, 5},
		{4, 1, -1},
		{5, 0, -1},
		{5, 6, -1},
	}
	for _, tt := range tests {
		result := getExpectedGroupSize(tt.players, tt.roundCount)
		if result != tt.expected {
			t.Errorf("getExpectedGroupSize(%d, %d) = %d, expected %d", tt.players, tt.roundCount, result, tt.expected)
		}
	}
}

func TestItoa(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{2, "2"},
		{3, "3"},
		{5, "5"},
		{10, "10"},
	}
	for _, tt := range tests {
		result := itoa(tt.input)
		if result != tt.expected {
			t.Errorf("itoa(%d) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestPlayerExists(t *testing.T) {
	g := types.NewGame("Test", "Owner1", "")
	g.Players = append(g.Players, "Owner1")
	if !playerExists(*g, "Owner1") {
		t.Error("expected Owner1 to exist")
	}
	if playerExists(*g, "NonExistent") {
		t.Error("expected NonExistent to not exist")
	}
}

func TestVerifyPlayerSelection(t *testing.T) {
	g := types.NewGame("Test", "Owner1", "")
	g.Players = append(g.Players, "Owner1", "P2", "P3")

	if !verifyPlayerSelection(*g, []string{"Owner1", "P2"}) {
		t.Error("expected valid selection")
	}
	if verifyPlayerSelection(*g, []string{"Owner1", "NonExistent"}) {
		t.Error("expected invalid selection for non-existent player")
	}
}
