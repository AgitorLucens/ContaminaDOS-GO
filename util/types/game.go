package types

import (
	"github.com/google/uuid"
)

type Game struct {
	Name   		 string     `json:"name" example:"Epsilon Centauri"`
	Owner    	 string     `json:"owner" example:"Thanos"`
	Status 		 GameStatus `json:"status"`
	PWD    		 string     `json:"-"`
	// Password is optional and should be omitted if empty
	Password     bool      `json:"password,omitempty" example:"false"`
	GameId       uuid.UUID `json:"id,omitempty"`
	CreatedAt    string    `json:"createdAt" example:"2023-10-01T12:00:00Z"`
	UpdatedAt    string    `json:"updatedAt" example:"2023-10-01T12:00:00Z"`
	Players      []string  `json:"players" example:"Thanos, Gamora"`
	Enemies      []string  `json:"enemies" example:"Star-Lord, Drax"`
	CurrentRound string    `json:"currentRound" example:"1"`
}

type GameStatus string

const (
	StatusLobby  GameStatus = "lobby"
	StatusRounds GameStatus = "rounds"
	StatusEnded  GameStatus = "ended"
)

func NewGame(name string, owner string, password string) *Game {
	return &Game{
		Name:    name,
		Owner:   owner,
		PWD:     password,
		Players: []string{},
		Enemies: []string{},
		
	}
}
func (g *Game) GetName() string {
	return g.Name
}
func (g *Game) SetName(name string) {
	g.Name = name
}
func (g *Game) GetOwner() string {
	return g.Owner
}
func (g *Game) SetOwner(owner string) {
	g.Owner = owner
}
func (g *Game) GetPassword() string {
	return g.PWD
}
func (g *Game) SetPassword(password string) {
	g.PWD = password
}

func (g *Game) GetGameId() uuid.UUID {
	return g.GameId
}
func (g *Game) SetGameId(gameId uuid.UUID) {
	g.GameId = gameId
}
func (g *Game) GetStatus() GameStatus {
	return g.Status
}
func (g *Game) SetStatus(status GameStatus) {
	g.Status = status
}
func (g *Game) GetCreatedAt() string {
	return g.CreatedAt
}
func (g *Game) SetCreatedAt(createdAt string) {
	g.CreatedAt = createdAt
}
func (g *Game) GetUpdatedAt() string {
	return g.UpdatedAt
}
func (g *Game) SetUpdatedAt(updatedAt string) {
	g.UpdatedAt = updatedAt
}
func (g *Game) GetRound() string {
	return g.CurrentRound
}
func (g *Game) SetRound(rounds string) {
	g.CurrentRound = rounds
}
func (g *Game) GetGame() Game {

	return *g
}

/*
Add players to the game
*/
func (g *Game) SetPlayers(player string) {
	g.Players = append(g.Players, player)
}

func (g Game) Equal(other Game) bool {
	return g.Password == other.Password &&
		g.CurrentRound == other.CurrentRound &&
		g.Status == other.Status
}
