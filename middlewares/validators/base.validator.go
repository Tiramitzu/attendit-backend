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

func PathUserIdValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("userId")
		err := validation.Validate(id, is.MongoID)
		if err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, "Invalid userId: "+id)
			return
		}

		c.Next()
	}
}

func PathAttendanceIdValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("attendanceId")
		err := validation.Validate(id, is.MongoID)
		if err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, "Invalid attendanceId: "+id)
			return
		}

		c.Next()
	}
}

func PathScheduleIdValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("scheduleId")
		err := validation.Validate(id, is.MongoID)
		if err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, "Invalid scheduleId: "+id)
			return
		}

		c.Next()
	}
}

func PathPageValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		page := c.Param("page")
		err := validation.Validate(page, is.Digit)
		if err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, "Invalid page: "+page)
			return
		}

		c.Next()
	}
}
