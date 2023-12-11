package validators

import (
	"net/http"

	"attendit/backend/models"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func PathIdValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		err := validation.Validate(id, is.MongoID)
		if err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, "invalid id: "+id)
			return
		}

		c.Next()
	}
}

func PathPageValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		page := c.Param("page")
		err := validation.Validate(page, validation.Required, validation.Min(1))
		if err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, "invalid page: "+page)
			return
		}

		c.Next()
	}
}
