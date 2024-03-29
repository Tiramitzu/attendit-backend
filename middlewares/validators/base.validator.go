package validators

import (
	"attendit/backend/models"
	"attendit/backend/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

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
		if err != nil && id != "@me" {
			models.SendErrorResponse(c, http.StatusBadRequest, "Invalid userId: "+id)
			return
		} else if id == "@me" {
			user, err := services.GetUserByToken(c.GetHeader("Authorization")[7:])
			if err != nil {
				models.SendErrorResponse(c, http.StatusUnauthorized, "Invalid token")
				return
			}

			c.Set("userId", user.ID.Hex())
			c.Next()
		} else {
			authorization := c.GetHeader("Authorization")
			if authorization == "" {
				models.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
				return
			}

			user, err := services.GetUserByToken(authorization[7:])
			if err != nil {
				models.SendErrorResponse(c, http.StatusUnauthorized, "Invalid token")
				return
			}

			userId, _ := primitive.ObjectIDFromHex(id)

			if user.AccessLevel < 1 && user.ID != userId {
				models.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
				return
			}

			c.Set("userId", id)
			c.Next()
		}
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
