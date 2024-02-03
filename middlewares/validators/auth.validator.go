package validators

import (
	"attendit/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func LoginValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var loginRequest models.LoginRequest
		_ = c.ShouldBindBodyWith(&loginRequest, binding.JSON)

		if err := loginRequest.Validate(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}

		c.Next()
	}
}
