package command

import (
	"be/internal/logic"
	"be/util/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Join Game
// @Description Gets the content of a game. Returns an error if it doesn't exist.
// @Tags Player
// @Accept json
// @Param gameId path string true "Game ID"
// @Param password header string true "Password"
// @Param player header string true "Player ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/games/{gameId} [put]
func JoinGame(c *gin.Context) {
	var req types.JoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var e types.Others
	//var games types.Game
	c.Param("gameId")
	g,err := logic.NewGamePool().JoinGame(req.Player, req)

	if err != nil {
		if err.Error() == "" {
			e.AppendError(types.Error{
				Status:  404,
				Message: "Contraseña Incorrecta",
			})
		}
		e.AppendError(types.Error{
			Status:  404,
			Message: err.Error(),
		})
	}
	if e.Others != nil {
		c.JSON(http.StatusCreated, types.Response{
			Status:  http.StatusBadRequest,
			Message: e.Others[0].Message,
			Data:    []types.Game{},
			Others:  e,
		})
		return
	}
	// Do Validations
	// if req.GameId == "" {
	c.JSON(http.StatusCreated, types.Response{
		Status:  http.StatusCreated,
		Message: "Game joined successfully",
		Data: g,
		Others: []types.Error{},
	})
}

// @Summary Game Start
// @Description Initiate Game if the amount of players is reach
// @Tags Player
// @Accept json
// @Param gameId path string true "Game ID"
// @Param password header string true "Password"
// @Param player header string true "Player ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/games/{gameId}/start [head]
func GameStart(c *gin.Context) {
	id := c.Param("gameId")

	if len(id) == 0 {
		c.Writer.Header().Set("x-msg","Missing header");
		c.Writer.Header().Set("status","400  Bad Request");
		return
	}
	req := types.StartGameRequest{
		GameId: id,
		Player: c.GetHeader("Player"),
		Password: c.GetHeader("Password"),
	}
	_, err := logic.NewGamePool().StartGame(req)
	if err != nil {
		c.Writer.Header().Set("x-msg", err.Error())
		c.Writer.Header().Set("status", "409 Conflict")
		return
	}
	c.Writer.Header().Set("x-msg","Started successfuly");
	c.Writer.Header().Set("status","200 OK");
}

func ProposeGroup(c *gin.Context) {
	gameId := c.Param("gameId")
	roundId := c.Param("roundId")
	player := c.GetHeader("player")

	var req types.ProposeGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid or missing group",
			Data:    nil,
			Others:  []types.Error{},
		})
		return
	}

	round, err := logic.NewGamePool().ProposeGroup(gameId, roundId, player, req.Group)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Round not found" || err.Error() == "Game not found" {
			status = http.StatusNotFound
		} else if err.Error() == "It is not the time for proposing groups" {
			status = http.StatusPreconditionRequired
		} else if err.Error() == "You are not the Leader" {
			status = http.StatusNotFound
		}
		c.JSON(status, types.Response{
			Status:  status,
			Message: err.Error(),
			Data:    nil,
			Others:  []types.Error{},
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Status:  http.StatusOK,
		Message: "Group Created",
		Data:    round,
		Others:  []types.Error{},
	})
}

func VoteGroup(c *gin.Context) {
	gameId := c.Param("gameId")
	roundId := c.Param("roundId")
	player := c.GetHeader("player")

	var req types.VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid or missing vote",
			Data:    nil,
			Others:  []types.Error{},
		})
		return
	}

	round, err := logic.NewGamePool().VoteGroup(gameId, roundId, player, req.Vote)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Round not found" || err.Error() == "Game not found" {
			status = http.StatusNotFound
		} else if err.Error() == "You have already voted" {
			status = http.StatusConflict
		} else if err.Error() == "It is not the time for voting" {
			status = http.StatusPreconditionRequired
		} else if err.Error() == "Player is not part of the game" {
			status = http.StatusBadRequest
		}
		c.JSON(status, types.Response{
			Status:  status,
			Message: err.Error(),
			Data:    nil,
			Others:  []types.Error{},
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Status:  http.StatusOK,
		Message: "Vote registered",
		Data:    round,
		Others:  []types.Error{},
	})
}

func SubmitAction(c *gin.Context) {
	gameId := c.Param("gameId")
	roundId := c.Param("roundId")
	player := c.GetHeader("player")

	var req types.ActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid or missing action",
			Data:    nil,
			Others:  []types.Error{},
		})
		return
	}

	round, err := logic.NewGamePool().SubmitAction(gameId, roundId, player, req.Action)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Round not found" || err.Error() == "Game not found" {
			status = http.StatusNotFound
		} else if err.Error() == "You cannot contribute in this round" {
			status = http.StatusForbidden
		} else if err.Error() == "It is not the time for actions" {
			status = http.StatusPreconditionRequired
		}
		c.JSON(status, types.Response{
			Status:  status,
			Message: err.Error(),
			Data:    nil,
			Others:  []types.Error{},
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Status:  http.StatusOK,
		Message: "Action registered",
		Data:    round,
		Others:  []types.Error{},
	})
}
