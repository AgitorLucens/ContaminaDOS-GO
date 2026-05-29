package querie


import (
	//"fmt"
	"github.com/gin-gonic/gin"	
)



// @Summary Get Hello Message
// @Description Returns a simple hello message
// @Tags Example
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/hello [get]
func GetHello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello, world!"})
}