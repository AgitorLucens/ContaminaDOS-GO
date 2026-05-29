package querie

import (
	"be/internal/logic"
	"be/util/types"
	"be/util/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Game Search
// @Description  passing in the appropriate options, you can search for games in the system
// @Tags Public, Player
// @Produce json
// @Param name query string false "gameId"
// @Param stauts query string false "password"
// @Param page header integer false "0"
// @Param limit query integer false "50"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/games [get]
func GameSearch(c *gin.Context) {
	// Get query parameters
	name := c.Query("name")
	status := c.Query("status")
	// Defaults
    page := c.DefaultQuery("page", "0")
    limit := c.DefaultQuery("limit", "50")
	
	httpStatus := http.StatusOK

	// Validate query parameters
	others := types.Others{}
	e := others.GetErrors(name, status, page, limit)
	if len(e) > 0 {
		httpStatus = http.StatusBadRequest
		c.JSON(httpStatus, types.Response{
			Status:  httpStatus,
			Message: "Bad Request",
			Data:    nil,
			Others:  types.Others{
				Others: e,
			},
		})
		return
	}

	p, err := utils.ToInt(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, e)
		return
	}
	l,err :=  utils.ToInt(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}
	games := logic.NewGamePool().SearchGamesByNameAndStatus(name, status, p, l)
	// Process the request and return a response

	var msg string
	if len(games) < 1 {
		msg = "Not Games Found"
	} else {
		msg = "Games Found"
	}

	c.JSON(200, &types.Response{
		Status:  httpStatus,
		Message: msg,
		Data:    games,
		Others:  e,
	})
}
