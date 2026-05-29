package command

import (
	"be/internal/logic"
	"be/util/types"
	"net/http"

	"github.com/gin-gonic/gin"
	//"fmt"
)

// @Summary Create Game
// @Description By passing in the appropriate options, you can create a game. Password is optional.
// @Tags Public, Player
// @Accept json
// @Produce json
// @Param game body types.CreateRequest true "Game object that needs to be created"
// @Success 200 {object} map[string]interface{} "Game Created"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Router /api/games [post]
func CreateGame(c *gin.Context) {
	var newGame types.CreateRequest

	// Bind JSON to newGame struct
	if err := c.ShouldBindJSON(&newGame); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Header or Body"})
		return
	}

	// Basic validation
	if newGame.Name == "" || newGame.Owner == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "400", "msg": "missing data"})
		return
	}

	// Create a new game
	g,err := logic.NewGamePool().CreateGame(*types.NewGame(newGame.Name, newGame.Owner, newGame.Password))

	if err != nil {
		if err.Error() == "Game already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"Status": 409, "msg": err.Error()})
			return
		}
	}
	var a interface{}
	// Respond success
	c.JSON(http.StatusOK, types.Response{
		Status:    http.StatusOK,
		Message: "Game created successfully",
		Data: &types.Game{
				Name: g.Name,
				PWD: g.PWD,
				Owner: g.Owner,
				GameId: g.GameId,
				CreatedAt: g.CreatedAt,
				UpdatedAt: g.UpdatedAt,
				Players: g.Players,
				Enemies: g.Enemies,
				CurrentRound: g.CurrentRound,
				Status: g.Status,
				},
		Others: a,
	})
}
