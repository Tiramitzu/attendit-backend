package validators

import (
	"attendit/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func FeedbackValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var feedbackRequest models.FeedbackRequest
		_ = c.ShouldBindBodyWith(&feedbackRequest, binding.JSON)

		if err := feedbackRequest.Validate(); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
		}

		c.Next()
	}
}
