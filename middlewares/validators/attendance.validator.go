package validators

import (
	"attendit/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CheckInValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var checkInRequest models.CheckInRequest
		_ = c.ShouldBindBodyWith(&checkInRequest, binding.JSON)

		if err := checkInRequest.Validate(); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
		}

		c.Next()
	}
}
