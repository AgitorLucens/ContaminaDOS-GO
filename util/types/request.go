package types

/*
Public Requests Objects
*/
type CreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password"`
	Owner    string `json:"owner" binding:"required"`
}

/*
Player Requests Objects
*/
type GetGameRequest struct {
	Password string `header:"Password"`
	Player   string `header:"Player" binding:"required"`
}

type JoinRequest struct {
	GameId   string `json:"gameId" binding:"required"`
	Password string `json:"password"`
	Player   string `json:"player" binding:"required"`
}

type StartGameRequest struct {
	GameId string
	Password string
	Player string
}

type RoundRequest struct {
	GameId   string
	RoundId  string
	Password string
	Player   string
}

type ProposeGroupRequest struct {
	Group []string `json:"group" binding:"required"`
}

type VoteRequest struct {
	Vote bool `json:"vote"`
}

type ActionRequest struct {
	Action bool `json:"action"`
}
