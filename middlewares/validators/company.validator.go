package validators

import (
	"attendit/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func CreateCompanyValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var createCompanyRequest models.CreateCompanyRequest
		_ = c.ShouldBindBodyWith(&createCompanyRequest, binding.JSON)

		if err := createCompanyRequest.Validate(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}

		c.Next()
	}
}

func InsertMembersToCompanyValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var insertMembersToCompanyRequest models.InsertMembersToCompanyRequest
		_ = c.ShouldBindBodyWith(&insertMembersToCompanyRequest, binding.JSON)

		if err := insertMembersToCompanyRequest.Validate(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		}

		c.Next()
	}
}

func CreateInvitationValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var createInvitationRequest models.CreateInvitationRequest
		_ = c.ShouldBindBodyWith(&createInvitationRequest, binding.JSON)

		if err := createInvitationRequest.Validate(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}

		c.Next()
	}
}
