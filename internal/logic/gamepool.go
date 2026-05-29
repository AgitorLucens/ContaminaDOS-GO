package logic

import (
	types "be/util/types"
	"errors"

	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GamePool struct {
	mutex      sync.Mutex
	games      map[string]types.Game
	gamesUUID  map[uuid.UUID]types.Game
	rounds     map[uuid.UUID]types.Round
	roundsUUID map[uuid.UUID]types.Round
}

var gamePoolInstance *GamePool
var once sync.Once

func NewGamePool() *GamePool {
	once.Do(func() {
		gamePoolInstance = &GamePool{
			games:      make(map[string]types.Game),
			gamesUUID:  make(map[uuid.UUID]types.Game),
			rounds:     make(map[uuid.UUID]types.Round),
			roundsUUID: make(map[uuid.UUID]types.Round),
		}
	})
	return gamePoolInstance
}

func (s *GamePool) CreateGame(game types.Game) (types.Game, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.games[game.Name]; exists {
		return types.Game{}, errors.New("Game already exists")
	}
	uId := uuid.New()
	game.SetGameId(uId)
	game.SetStatus(types.StatusLobby)
	game.SetRound("0000000000000000000000000")
	game.SetCreatedAt(string(time.Now().Format("2006-01-02 15:04:05")))
	game.Players = append(game.Players, game.Owner)
	game.Enemies = []string{}
	s.games[game.Name] = game
	s.gamesUUID[uId] = game
	return game, nil
}

func (s *GamePool) SearchGamesByNameAndStatus(name, status string, page, limit int) []types.Game {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var result []types.Game
	for _, game := range s.games {
		if name != "" && !strings.Contains(strings.ToLower(game.Name), strings.ToLower(name)) {
			continue
		}
		if status != "" && status != string(game.Status) {
			continue
		}
		result = append(result, game)
	}

	start := page * limit
	if start >= len(result) {
		return []types.Game{}
	}

	end := start + limit
	if end > len(result) {
		end = len(result)
	}

	return result[start:end]
}

func (s *GamePool) GetGame(Id string, r types.GetGameRequest) (*types.Game, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	gameId, err := uuid.Parse(Id)
	if err != nil {
		return nil, errors.New("Invalid game ID")
	}

	game, exists := s.gamesUUID[gameId]
	if !exists {
		return nil, errors.New("Game not found")
	}

	return &game, nil
}

func (s *GamePool) JoinGame(playername string, game types.JoinRequest) (*types.Game, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	str := game.GameId
	id, err := uuid.Parse(str)
	if err != nil {
		return nil, errors.New("Invalid game ID")
	}
	g, exists := s.gamesUUID[id]
	if !exists {
		return nil, errors.New("Game not found")
	}
	if g.PWD != game.Password {
		return nil, errors.New("Invalid game password")
	}

	err = AddPlayer(playername, &g)

	if err != nil {
		return nil, err
	}
	s.gamesUUID[id] = g
	s.games[g.Name] = g
	return &g, nil
}

func (s *GamePool) StartGame(game types.StartGameRequest) (*types.Game, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	str := game.GameId
	id, err := uuid.Parse(str)
	if err != nil {
		return nil, errors.New("Invalid game ID")
	}
	g, exists := s.gamesUUID[id]
	if !exists {
		return nil, errors.New("Game not found")
	}
	if g.PWD != game.Password {
		return nil, errors.New("Invalid game password")
	}

	if g.GetStatus() != types.StatusLobby {
		return nil, errors.New("Game already started")
	}

	if len(g.Players) < 5 {
		return nil, errors.New("Need 5 players to start")
	}

	if g.Owner != game.Player {
		return nil, errors.New("You are not the owner of this game")
	}

	enemiesCount := getEnemiesCountAtStart(len(g.Players))
	availablePlayers := make([]string, len(g.Players))
	copy(availablePlayers, g.Players)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(availablePlayers), func(i, j int) {
		availablePlayers[i], availablePlayers[j] = availablePlayers[j], availablePlayers[i]
	})

	g.Enemies = availablePlayers[:enemiesCount]
	

	r := types.NewRound(g.GameId)
	r.SetLeader(getRandomLeader(g.Players))
	s.roundsUUID[r.Id] = *r
	s.rounds[r.GameId] = *r

	g.SetRound(r.Id.String())
	g.SetStatus(types.StatusRounds)

	s.gamesUUID[id] = g
	s.games[g.Name] = g

	return &g, nil
}

func (s *GamePool) GetRounds(gameId string) ([]types.Round, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	gId, err := uuid.Parse(gameId)
	if err != nil {
		return nil, errors.New("Invalid game ID")
	}

	var result []types.Round
	for _, r := range s.roundsUUID {
		if r.GameId == gId {
			result = append(result, r)
		}
	}

	return result, nil
}

func (s *GamePool) ShowRound(gameId string, roundId string) (*types.Round, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	gId, err := uuid.Parse(gameId)
	if err != nil {
		return nil, errors.New("Invalid game ID")
	}

	rId, err := uuid.Parse(roundId)
	if err != nil {
		return nil, errors.New("Invalid round ID")
	}

	round, exists := s.roundsUUID[rId]
	if !exists || round.GameId != gId {
		return nil, errors.New("Round not found")
	}

	return &round, nil
}

func (s *GamePool) ProposeGroup(gameId string, roundId string, player string, group []string) (*types.Round, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	gId, err := uuid.Parse(gameId)
	if err != nil {
		return nil, errors.New("Invalid game ID")
	}

	rId, err := uuid.Parse(roundId)
	if err != nil {
		return nil, errors.New("Invalid round ID")
	}

	g, exists := s.gamesUUID[gId]
	if !exists {
		return nil, errors.New("Game not found")
	}

	r, exists := s.roundsUUID[rId]
	if !exists || r.GameId != gId {
		return nil, errors.New("Round not found")
	}

	if r.Leader != player {
		return nil, errors.New("You are not the Leader")
	}

	if r.Status != types.RoundStatusWaitingOnLeader {
		return nil, errors.New("It is not the time for proposing groups")
	}

	if !verifyPlayerSelection(g, group) {
		return nil, errors.New("Group member is not part of the game")
	}

	if !verifyGroupCount(g, s.getRoundsForGame(gId), group) {
		return nil, errors.New(getResponseGroupCount(g, s.getRoundsForGame(gId)))
	}

	r.Group = group
	r.Status = types.RoundStatusVoting
	s.roundsUUID[rId] = r
	s.rounds[r.GameId] = r

	return &r, nil
}

func (s *GamePool) VoteGroup(gameId string, roundId string, player string, vote bool) (*types.Round, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	gId, err := uuid.Parse(gameId)
	if err != nil {
		return nil, errors.New("Invalid game ID")
	}

	rId, err := uuid.Parse(roundId)
	if err != nil {
		return nil, errors.New("Invalid round ID")
	}

	g, exists := s.gamesUUID[gId]
	if !exists {
		return nil, errors.New("Game not found")
	}

	r, exists := s.roundsUUID[rId]
	if !exists || r.GameId != gId {
		return nil, errors.New("Round not found")
	}

	if !playerExists(g, player) {
		return nil, errors.New("Player is not part of the game")
	}

	if r.Status != types.RoundStatusVoting {
		return nil, errors.New("It is not the time for voting")
	}

	if r.HasVoted(player) {
		return nil, errors.New("You have already voted")
	}

	r.Votes = append(r.Votes, vote)
	r.AlreadyVote = append(r.AlreadyVote, player)

	if len(r.Votes) == len(g.Players) {
		trueCount := 0
		falseCount := 0
		for _, v := range r.Votes {
			if v {
				trueCount++
			} else {
				falseCount++
			}
		}

		if trueCount > falseCount {
			r.Status = types.RoundStatusWaitingOnGroup
		} else if trueCount < falseCount && r.Phase != types.RoundPhaseVote3 {
			r.ClearGroupAndVotes()
			r.Status = types.RoundStatusWaitingOnLeader
			r.AdvancePhase()
		} else if trueCount < falseCount && r.Phase == types.RoundPhaseVote3 {
			r.Result = types.RoundResultEnemies
			r.Status = types.RoundStatusEnded

			if verifyGameWinner(s, gId, g) {
				g.SetStatus(types.StatusEnded)
				s.gamesUUID[gId] = g
				s.games[g.Name] = g
			} else {
				newRound := types.NewRound(g.GameId)
				newRound.SetLeader(getRandomLeader(g.Players))
				s.roundsUUID[newRound.Id] = *newRound
				s.rounds[newRound.GameId] = *newRound
				g.SetRound(newRound.Id.String())
				s.gamesUUID[gId] = g
				s.games[g.Name] = g
			}
		}
	}

	s.roundsUUID[rId] = r
	s.rounds[r.GameId] = r

	return &r, nil
}

func (s *GamePool) SubmitAction(gameId string, roundId string, player string, action bool) (*types.Round, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	gId, err := uuid.Parse(gameId)
	if err != nil {
		return nil, errors.New("Invalid game ID")
	}

	rId, err := uuid.Parse(roundId)
	if err != nil {
		return nil, errors.New("Invalid round ID")
	}

	g, exists := s.gamesUUID[gId]
	if !exists {
		return nil, errors.New("Game not found")
	}

	r, exists := s.roundsUUID[rId]
	if !exists || r.GameId != gId {
		return nil, errors.New("Round not found")
	}

	if !playerExists(g, player) {
		return nil, errors.New("Player is not part of the game")
	}

	if !r.IsPlayerInGroup(player) {
		return nil, errors.New("You cannot contribute in this round")
	}

	if r.Status != types.RoundStatusWaitingOnGroup {
		return nil, errors.New("It is not the time for actions")
	}

	r.Actions = append(r.Actions, action)

	if len(r.Actions) == len(r.Group) {
		hasSabotage := false
		for _, a := range r.Actions {
			if !a {
				hasSabotage = true
				break
			}
		}

		if hasSabotage {
			r.Result = types.RoundResultEnemies
		} else {
			r.Result = types.RoundResultCitizens
		}

		if verifyGameWinner(s, gId, g) {
			g.SetStatus(types.StatusEnded)
			r.Status = types.RoundStatusEnded
			s.gamesUUID[gId] = g
			s.games[g.Name] = g
		} else {
			r.Status = types.RoundStatusEnded
			newRound := types.NewRound(g.GameId)
			newRound.SetLeader(getRandomLeader(g.Players))
			s.roundsUUID[newRound.Id] = *newRound
			s.rounds[newRound.GameId] = *newRound
			g.SetRound(newRound.Id.String())
			s.gamesUUID[gId] = g
			s.games[g.Name] = g
		}
	}

	s.roundsUUID[rId] = r
	s.rounds[r.GameId] = r

	return &r, nil
}

func AddPlayer(pname string, g *types.Game) error {
	for i := 0; len(g.Players) > i; i++ {
		if g.Players[i] == pname {
			return errors.New("Jugador ya existe")
		}
	}
	g.Players = append(g.Players, pname)
	return nil
}

func playerExists(g types.Game, player string) bool {
	for _, p := range g.Players {
		if p == player {
			return true
		}
	}
	return false
}

func getRandomLeader(players []string) string {
	rand.Seed(time.Now().UnixNano())
	return players[rand.Intn(len(players))]
}

func getEnemiesCountAtStart(playerCount int) int {
	if playerCount == 5 || playerCount == 6 {
		return 2
	} else if playerCount >= 7 && playerCount <= 9 {
		return 3
	}
	return 4
}

func (s *GamePool) getRoundsForGame(gameId uuid.UUID) []types.Round {
	var result []types.Round
	for _, r := range s.roundsUUID {
		if r.GameId == gameId {
			result = append(result, r)
		}
	}
	return result
}

func verifyPlayerSelection(g types.Game, group []string) bool {
	for _, member := range group {
		if !playerExists(g, member) {
			return false
		}
	}
	return true
}

func verifyGroupCount(g types.Game, rounds []types.Round, group []string) bool {
	playerCount := len(g.Players)
	roundCount := len(rounds)
	groupCount := len(group)

	expected := getExpectedGroupSize(playerCount, roundCount)
	return groupCount == expected
}

func getExpectedGroupSize(playerCount, roundCount int) int {
	table := map[int][]int{
		5:  {2, 3, 2, 3, 3},
		6:  {2, 3, 4, 3, 4},
		7:  {2, 3, 3, 4, 4},
		8:  {3, 4, 4, 5, 5},
		9:  {3, 4, 4, 5, 5},
		10: {3, 4, 4, 5, 5},
	}

	sizes, ok := table[playerCount]
	if !ok || roundCount < 1 || roundCount > len(sizes) {
		return -1
	}
	return sizes[roundCount-1]
}

func getResponseGroupCount(g types.Game, rounds []types.Round) string {
	expected := getExpectedGroupSize(len(g.Players), len(rounds))
	if expected < 0 {
		return "Invalid group configuration"
	}
	return "Requires a group of " + itoa(expected) + " members"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

func verifyGameWinner(s *GamePool, gameId uuid.UUID, g types.Game) bool {
	citizensCount := 0
	enemiesCount := 0
	for _, r := range s.roundsUUID {
		if r.GameId == gameId {
			if r.Result == types.RoundResultCitizens {
				citizensCount++
			} else if r.Result == types.RoundResultEnemies {
				enemiesCount++
			}
		}
	}
	return citizensCount == 3 || enemiesCount == 3
}
