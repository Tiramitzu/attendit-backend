package middlewares

import (
	"attendit/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppRecovery() func(c *gin.Context, recovered interface{}) {
	return func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			models.SendErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false}) // recovery failed
	}
}
