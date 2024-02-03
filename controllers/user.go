package controllers

import (
	"attendit/backend/models"
	"attendit/backend/services"
	redisServices "attendit/backend/services/redis"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

// GetCurrentUser godoc
// @Summary      GetCurrentUser
// @Description  gets the current user
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{userId} [get]
func GetCurrentUser(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userId, _ := c.Get("userId")

	user, err := redisServices.GetUserFromCache(userId.(primitive.ObjectID))
	if err == nil {
		response.StatusCode = http.StatusOK
		response.Success = true
		response.Data = gin.H{"user": user, "cache": true}
		response.SendResponse(c)
		return
	}

	user, err = services.GetUserByToken(c.GetHeader("Authorization")[7:])
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUser(user)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"user": user}
	response.SendResponse(c)
}

// ModifyCurrentUser godoc
// @Summary      ModifyCurrentUser
// @Description  modifies the current user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        req  body      models.ModifyUserRequest    true  "Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{userId} [patch]
func ModifyCurrentUser(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var requestBody models.ModifyUserRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	user, err := services.GetUserByToken(c.GetHeader("Authorization")[7:])
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	user.Email = requestBody.Email
	user.FullName = requestBody.FullName
	user.Phone = requestBody.Phone
	user.Password = requestBody.Password

	updatedUser, err := services.UpdateUser(user)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUser(updatedUser)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"user": updatedUser}
	response.SendResponse(c)
}

// GetUserAttendances godoc
// @Summary      GetUserAttendances
// @Description  gets the user attendances
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{userId}/attendances/:page [get]
func GetUserAttendances(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	page, _ := strconv.Atoi(c.Param("page"))
	if page == 0 {
		page = 1
	}

	user, err := services.GetUserByToken(c.GetHeader("Authorization")[7:])
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	// Try to get attendances from cache
	attendances, cacheErr := redisServices.GetUserAttendancesFromCache(user.ID, page)
	if cacheErr == nil {
		response.StatusCode = http.StatusOK
		response.Success = true
		response.Data = gin.H{"attendances": attendances, "cache": true}
		response.SendResponse(c)
		return
	}

	// If cache retrieval fails, get attendances from services
	attendances, err = services.GetUserAttendances(user.ID, page)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	// If attendances is nil, return a success response
	if attendances != nil {
		// Cache attendances for future use
		redisServices.CacheUserAttendancesByCompany(user.ID, attendances, page)
	}

	// Send a success response
	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"attendances": attendances}
	response.SendResponse(c)
}
