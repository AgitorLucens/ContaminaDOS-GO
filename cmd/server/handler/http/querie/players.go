package querie

import (
	"be/internal/logic"
	"be/util/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Game
// @Description Gets the content of a game. Returns an error if it doesn't exist.
// @Tags Player
// @Accept json
// @Param gameId path string true "Game ID"
// @Param password header string true "Password"
// @Param player header string true "Player ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/games/{gameId} [get]
func GetGame(c *gin.Context) {
	req := types.GetGameRequest{
		Password: c.GetHeader("Password"),
		Player:   c.GetHeader("Player"),
	}

	gameId := c.Param("gameId")
	g, err := logic.NewGamePool().GetGame(gameId, req)
	if err != nil {
		c.JSON(http.StatusOK, &types.Response{
			Status:  http.StatusOK,
			Message: "Game Not Found",
			Data:    g,
			Others:  []types.Error{},
		})
		return
	}

	c.JSON(http.StatusOK, &types.Response{
		Status:  http.StatusOK,
		Message: "Game Found",
		Data:    g,
		Others:  []types.Error{},
	})
}

// @Summary Get Rounds
// @Description Gets the rounds of a game.
// @Tags Player
// @Accept json
// @Param gameId path string true "Game ID"
// @Param password header string true "Password"
// @Param player header string true "Player ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/games/{gameId}/rounds [get]
func GetRounds(c *gin.Context) {
	gameId := c.Param("gameId")

	rounds, err := logic.NewGamePool().GetRounds(gameId)
	if err != nil {
		c.JSON(http.StatusNotFound, types.Response{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    []types.Round{},
			Others:  []types.Error{},
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Status:  http.StatusOK,
		Message: "Results found",
		Data:    rounds,
		Others:  []types.Error{},
	})
}

// @Summary Show Round
// @Description Gets the content of a specific round.
// @Tags Player
// @Accept json
// @Param gameId path string true "Game ID"
// @Param roundId path string true "Round ID"
// @Param password header string true "Password"
// @Param player header string true "Player ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/games/{gameId}/rounds/{roundId} [get]
func ShowRound(c *gin.Context) {
	gameId := c.Param("gameId")
	roundId := c.Param("roundId")

	round, err := logic.NewGamePool().ShowRound(gameId, roundId)
	if err != nil {
		c.JSON(http.StatusNotFound, types.Response{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
			Others:  []types.Error{},
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Status:  http.StatusOK,
		Message: "Results found",
		Data:    round,
		Others:  []types.Error{},
	})
}
