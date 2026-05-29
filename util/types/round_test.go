package types

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewRound(t *testing.T) {
	gameId := uuid.New()
	r := NewRound(gameId)
	if r.GameId != gameId {
		t.Errorf("expected gameId %v, got %v", gameId, r.GameId)
	}
	if r.Status != RoundStatusWaitingOnLeader {
		t.Errorf("expected status waiting-on-leader, got %s", r.Status)
	}
	if r.Result != RoundResultNone {
		t.Errorf("expected result none, got %s", r.Result)
	}
	if r.Phase != RoundPhaseVote1 {
		t.Errorf("expected phase vote1, got %s", r.Phase)
	}
	if len(r.Group) != 0 {
		t.Errorf("expected empty group, got %d", len(r.Group))
	}
	if len(r.Votes) != 0 {
		t.Errorf("expected empty votes, got %d", len(r.Votes))
	}
	if len(r.Actions) != 0 {
		t.Errorf("expected empty actions, got %d", len(r.Actions))
	}
	if len(r.AlreadyVote) != 0 {
		t.Errorf("expected empty alreadyVote, got %d", len(r.AlreadyVote))
	}
	if r.Id == uuid.Nil {
		t.Error("expected non-nil round ID")
	}
}

func TestRoundGettersSetters(t *testing.T) {
	t.Parallel()

	t.Run("Id", func(t *testing.T) {
		r := NewRound(uuid.New())
		id := uuid.New()
		r.SetId(id)
		if r.GetId() != id {
			t.Errorf("got %v; want %v", r.GetId(), id)
		}
	})

	t.Run("Status", func(t *testing.T) {
		r := NewRound(uuid.New())
		r.SetStatus(RoundStatusVoting)
		if r.GetStatus() != RoundStatusVoting {
			t.Errorf("got %q; want %q", r.GetStatus(), RoundStatusVoting)
		}
	})

	t.Run("Result", func(t *testing.T) {
		r := NewRound(uuid.New())
		r.SetResult(RoundResultCitizens)
		if r.GetResult() != RoundResultCitizens {
			t.Errorf("got %q; want %q", r.GetResult(), RoundResultCitizens)
		}
	})

	t.Run("Phase", func(t *testing.T) {
		r := NewRound(uuid.New())
		r.SetPhase(RoundPhaseVote2)
		if r.GetPhase() != RoundPhaseVote2 {
			t.Errorf("got %q; want %q", r.GetPhase(), RoundPhaseVote2)
		}
	})

	t.Run("Leader", func(t *testing.T) {
		r := NewRound(uuid.New())
		r.SetLeader("Player1")
		if r.GetLeader() != "Player1" {
			t.Errorf("got %q; want %q", r.GetLeader(), "Player1")
		}
	})

	t.Run("Group", func(t *testing.T) {
		r := NewRound(uuid.New())
		group := []string{"P1", "P2"}
		r.SetGroup(group)
		got := r.GetGroup()
		if len(got) != 2 || got[0] != "P1" || got[1] != "P2" {
			t.Errorf("got %v; want %v", got, group)
		}
	})

	t.Run("Votes", func(t *testing.T) {
		r := NewRound(uuid.New())
		votes := []bool{true, false}
		r.SetVotes(votes)
		got := r.GetVotes()
		if len(got) != 2 || got[0] != true || got[1] != false {
			t.Errorf("got %v; want %v", got, votes)
		}
	})

	t.Run("Actions", func(t *testing.T) {
		r := NewRound(uuid.New())
		actions := []bool{true}
		r.SetActions(actions)
		got := r.GetActions()
		if len(got) != 1 || got[0] != true {
			t.Errorf("got %v; want %v", got, actions)
		}
	})

	t.Run("AlreadyVote", func(t *testing.T) {
		r := NewRound(uuid.New())
		expected := []string{"P1", "P2"}
		r.SetAlreadyVote(expected)
		got := r.GetAlreadyVote()
		if len(got) != 2 || got[0] != "P1" || got[1] != "P2" {
			t.Errorf("got %v; want %v", got, expected)
		}
	})
}

func TestRound_HasVoted(t *testing.T) {
	r := NewRound(uuid.New())
	r.AlreadyVote = []string{"P1", "P2"}
	if !r.HasVoted("P1") {
		t.Error("expected P1 to have voted")
	}
	if r.HasVoted("P3") {
		t.Error("expected P3 to not have voted")
	}
}

func TestRound_IsPlayerInGroup(t *testing.T) {
	r := NewRound(uuid.New())
	r.Group = []string{"P1", "P2"}
	if !r.IsPlayerInGroup("P1") {
		t.Error("expected P1 to be in group")
	}
	if r.IsPlayerInGroup("P3") {
		t.Error("expected P3 to not be in group")
	}
}

func TestRound_AdvancePhase(t *testing.T) {
	t.Parallel()

	r := NewRound(uuid.New())

	steps := []struct {
		name     string
		advance  bool
		expected RoundPhase
	}{
		{"starts at vote1", false, RoundPhaseVote1},
		{"advance to vote2", true, RoundPhaseVote2},
		{"advance to vote3", true, RoundPhaseVote3},
		{"stay at vote3", true, RoundPhaseVote3},
	}

	for _, step := range steps {
		t.Run(step.name, func(t *testing.T) {
			if step.advance {
				r.AdvancePhase()
			}
			if r.Phase != step.expected {
				t.Errorf("phase = %q; want %q", r.Phase, step.expected)
			}
		})
	}
}

func TestRound_ClearGroupAndVotes(t *testing.T) {
	r := NewRound(uuid.New())
	r.Group = []string{"P1", "P2"}
	r.Votes = []bool{true, false}
	r.AlreadyVote = []string{"P1", "P2"}

	r.ClearGroupAndVotes()

	if len(r.Group) != 0 {
		t.Errorf("expected empty group, got %d", len(r.Group))
	}
	if len(r.Votes) != 0 {
		t.Errorf("expected empty votes, got %d", len(r.Votes))
	}
	if len(r.AlreadyVote) != 0 {
		t.Errorf("expected empty alreadyVote, got %d", len(r.AlreadyVote))
	}
}

func TestRound_GetGameId(t *testing.T) {
	gameId := uuid.New()
	r := NewRound(gameId)
	if r.GetGameId() != gameId {
		t.Errorf("expected gameId %v, got %v", gameId, r.GetGameId())
	}
}
