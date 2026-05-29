package types

import (
	"be/util/utils"
)

type Error struct {
	Status    int    `json:"status"`

	Message string `json:"msg"`
}

type Others struct {
	// Status    int         `json:"status"`
	Others []Error `json:"others"`
}

func (o *Others) GetErrors(name string, status string, page string, limit string) []Error {
	others := []Error{}
	
	n, err := utils.ToInt(limit)
	if err != nil {
		others = append(others, Error{Status: 400, Message: "Invalid gameId"})
	}
	if n < 0 {
		others = append(others, Error{Status: 400, Message: "gameId must be greater than 0"})
	}
	p, err := utils.ToInt(page)
	if err != nil {
		others = append(others, Error{Status: 400, Message: "Invalid page number"})
	}
	if p < 0 {
		others = append(others, Error{Status: 400, Message: "page less than 0"})
	}
	if status == string(StatusLobby) || status == string(StatusRounds) || status == string(StatusEnded) {	
		others = append(others, Error{Status: 400, Message: "Invalid game status"})
	}
	if page == "" {
		others = append(others, Error{Status: 400, Message: "Page is required"})
	}
	if limit == "" {
		others = append(others, Error{Status: 400, Message: "Limit is required"})
	}
	return others
}

func (o *Others) AppendError(e Error){
	o.Others = append(o.Others,e)
}