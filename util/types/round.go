package types

import "github.com/google/uuid"

type RoundStatus string

const (
	RoundStatusWaitingOnLeader RoundStatus = "waiting-on-leader"
	RoundStatusVoting          RoundStatus = "voting"
	RoundStatusWaitingOnGroup  RoundStatus = "waiting-on-group"
	RoundStatusEnded           RoundStatus = "ended"
)

type RoundResult string

const (
	RoundResultNone      RoundResult = "none"
	RoundResultCitizens  RoundResult = "citizens"
	RoundResultEnemies   RoundResult = "enemies"
)

type RoundPhase string

const (
	RoundPhaseVote1 RoundPhase = "vote1"
	RoundPhaseVote2 RoundPhase = "vote2"
	RoundPhaseVote3 RoundPhase = "vote3"
)

type Round struct {
	Id          uuid.UUID   `json:"id"`
	GameId      uuid.UUID   `json:"-"`
	Leader      string      `json:"leader"`
	Status      RoundStatus `json:"status"`
	Result      RoundResult `json:"result"`
	Phase       RoundPhase  `json:"phase"`
	Group       []string    `json:"group"`
	Votes       []bool      `json:"votes"`
	Actions     []bool      `json:"-"`
	AlreadyVote []string `json:"-"`
}

func NewRound(gameId uuid.UUID) *Round {
	return &Round{
		Id:          uuid.New(),
		GameId:      gameId,
		Status:      RoundStatusWaitingOnLeader,
		Result:      RoundResultNone,
		Phase:       RoundPhaseVote1,
		Group:       []string{},
		Votes:       []bool{},
		Actions:     []bool{},
		AlreadyVote: []string{},
	}
}

func (r *Round) GetId() uuid.UUID {
	return r.Id
}

func (r *Round) SetId(id uuid.UUID) {
	r.Id = id
}

func (r *Round) GetGameId() uuid.UUID {
	return r.GameId
}

func (r *Round) GetStatus() RoundStatus {
	return r.Status
}

func (r *Round) SetStatus(status RoundStatus) {
	r.Status = status
}

func (r *Round) GetResult() RoundResult {
	return r.Result
}

func (r *Round) SetResult(result RoundResult) {
	r.Result = result
}

func (r *Round) GetPhase() RoundPhase {
	return r.Phase
}

func (r *Round) SetPhase(phase RoundPhase) {
	r.Phase = phase
}

func (r *Round) GetLeader() string {
	return r.Leader
}

func (r *Round) SetLeader(leader string) {
	r.Leader = leader
}

func (r *Round) GetGroup() []string {
	return r.Group
}

func (r *Round) SetGroup(group []string) {
	r.Group = group
}

func (r *Round) GetVotes() []bool {
	return r.Votes
}

func (r *Round) SetVotes(votes []bool) {
	r.Votes = votes
}

func (r *Round) GetActions() []bool {
	return r.Actions
}

func (r *Round) SetActions(actions []bool) {
	r.Actions = actions
}

func (r *Round) GetAlreadyVote() []string {
	return r.AlreadyVote
}

func (r *Round) SetAlreadyVote(alreadyVote []string) {
	r.AlreadyVote = alreadyVote
}

func (r *Round) HasVoted(player string) bool {
	for _, p := range r.AlreadyVote {
		if p == player {
			return true
		}
	}
	return false
}

func (r *Round) IsPlayerInGroup(player string) bool {
	for _, p := range r.Group {
		if p == player {
			return true
		}
	}
	return false
}

func (r *Round) AdvancePhase() {
	if r.Phase == RoundPhaseVote1 {
		r.Phase = RoundPhaseVote2
	} else if r.Phase == RoundPhaseVote2 {
		r.Phase = RoundPhaseVote3
	} else {
		r.Phase = RoundPhaseVote3
	}
}

func (r *Round) ClearGroupAndVotes() {
	r.Group = []string{}
	r.Votes = []bool{}
	r.AlreadyVote = []string{}
}
