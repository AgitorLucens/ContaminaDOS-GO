package types

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewGame(t *testing.T) {
	g := NewGame("TestGame", "Owner1", "secret")
	if g.Name != "TestGame" {
		t.Errorf("expected name TestGame, got %s", g.Name)
	}
	if g.Owner != "Owner1" {
		t.Errorf("expected owner Owner1, got %s", g.Owner)
	}
	if g.PWD != "secret" {
		t.Errorf("expected PWD secret, got %s", g.PWD)
	}
	if len(g.Players) != 0 {
		t.Errorf("expected empty players, got %d", len(g.Players))
	}
}

func TestNewGame_NoPassword(t *testing.T) {
	g := NewGame("TestGame", "Owner1", "")
	if g.PWD != "" {
		t.Errorf("expected empty PWD, got %s", g.PWD)
	}
}

func TestGameGettersSetters(t *testing.T) {
	t.Parallel()

	t.Run("Name", func(t *testing.T) {
		g := NewGame("Old", "Owner", "")
		g.SetName("New")
		if g.GetName() != "New" {
			t.Errorf("got %q; want %q", g.GetName(), "New")
		}
	})

	t.Run("Owner", func(t *testing.T) {
		g := NewGame("Test", "OldOwner", "")
		g.SetOwner("NewOwner")
		if g.GetOwner() != "NewOwner" {
			t.Errorf("got %q; want %q", g.GetOwner(), "NewOwner")
		}
	})

	t.Run("Password", func(t *testing.T) {
		g := NewGame("Test", "Owner", "old")
		g.SetPassword("new")
		if g.GetPassword() != "new" {
			t.Errorf("got %q; want %q", g.GetPassword(), "new")
		}
	})

	t.Run("GameId", func(t *testing.T) {
		g := NewGame("Test", "Owner", "")
		id := uuid.New()
		g.SetGameId(id)
		if g.GetGameId() != id {
			t.Errorf("got %v; want %v", g.GetGameId(), id)
		}
	})

	t.Run("Status", func(t *testing.T) {
		g := NewGame("Test", "Owner", "")
		g.SetStatus(StatusRounds)
		if g.GetStatus() != StatusRounds {
			t.Errorf("got %q; want %q", g.GetStatus(), StatusRounds)
		}
	})

	t.Run("CreatedAt", func(t *testing.T) {
		g := NewGame("Test", "Owner", "")
		g.SetCreatedAt("2023-01-01")
		if g.GetCreatedAt() != "2023-01-01" {
			t.Errorf("got %q; want %q", g.GetCreatedAt(), "2023-01-01")
		}
	})

	t.Run("UpdatedAt", func(t *testing.T) {
		g := NewGame("Test", "Owner", "")
		g.SetUpdatedAt("2023-01-02")
		if g.GetUpdatedAt() != "2023-01-02" {
			t.Errorf("got %q; want %q", g.GetUpdatedAt(), "2023-01-02")
		}
	})

	t.Run("Round", func(t *testing.T) {
		g := NewGame("Test", "Owner", "")
		g.SetRound("round-id-123")
		if g.GetRound() != "round-id-123" {
			t.Errorf("got %q; want %q", g.GetRound(), "round-id-123")
		}
	})
}

func TestGame_SetPlayers(t *testing.T) {
	g := NewGame("Test", "Owner", "")
	g.SetPlayers("P1")
	g.SetPlayers("P2")
	if len(g.Players) != 2 {
		t.Errorf("expected 2 players, got %d", len(g.Players))
	}
	if g.Players[0] != "P1" || g.Players[1] != "P2" {
		t.Errorf("expected [P1, P2], got %v", g.Players)
	}
}

func TestGame_GetGame(t *testing.T) {
	g := NewGame("Test", "Owner", "")
	g.SetStatus(StatusLobby)
	copy := g.GetGame()
	if copy.Name != "Test" {
		t.Errorf("expected name Test, got %s", copy.Name)
	}
	if copy.Status != StatusLobby {
		t.Errorf("expected status lobby, got %s", copy.Status)
	}
}

func TestGame_Equal(t *testing.T) {
	g1 := NewGame("Test", "Owner", "pass")
	g1.SetStatus(StatusLobby)
	g1.SetRound("round1")

	g2 := NewGame("Test", "Owner", "pass")
	g2.SetStatus(StatusLobby)
	g2.SetRound("round1")

	if !g1.Equal(*g2) {
		t.Error("expected games to be equal")
	}

	g2.SetStatus(StatusRounds)
	if g1.Equal(*g2) {
		t.Error("expected games to not be equal")
	}
}

func TestGameStatus_Constants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		got      GameStatus
		expected GameStatus
	}{
		{StatusLobby, "lobby"},
		{StatusRounds, "rounds"},
		{StatusEnded, "ended"},
	}

	for _, tt := range tests {
		if tt.got != tt.expected {
			t.Errorf("got %q; want %q", tt.got, tt.expected)
		}
	}
}
