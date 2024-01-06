package controllers

import (
	"attendit/backend/models"
	db "attendit/backend/models/db"
	"attendit/backend/services"
	redisServices "attendit/backend/services/redis"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

// GetUserSchedules godoc
// @Summary Get user schedules
// @Description Get user schedules
// @Tags schedule
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param page path int false "Page number"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /user/{userId}/schedules/:page [get]
func GetUserSchedules(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	page, _ := strconv.Atoi(c.Param("page"))
	if page == 0 {
		page = 1
	}

	userId, _ := c.Get("userId")
	user, _ := services.GetUserById(userId.(primitive.ObjectID))

	schedules, err := redisServices.GetUserSchedulesFromCache(user.ID, page)
	if err == nil {
		response.StatusCode = http.StatusOK
		response.Success = true
		response.Data = gin.H{"schedules": schedules, "cache": true}
		response.SendResponse(c)
		return
	}

	schedules, err = services.GetUserSchedules(user.ID, page)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUserSchedules(user.ID, schedules, page)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"schedules": schedules}
	response.SendResponse(c)
}

// GetUserSchedule godoc
// @Summary Get user schedule
// @Description Get user schedule
// @Tags schedule
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param scheduleId path string true "Schedule ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /user/{userId}/schedules/{scheduleId} [get]
// CreateUserSchedule godoc
// @Summary Create user schedule
// @Description Create user schedule
// @Tags schedule
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param title body string true "Schedule title"
// @Param startTime body string true "Schedule start time"
// @Param endTime body string true "Schedule end time"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /user/{userId}/schedules [post]
func CreateUserSchedule(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var requestBody models.ScheduleRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	userId, _ := c.Get("userId")
	user, _ := services.GetUserById(userId.(primitive.ObjectID))

	schedule := db.NewSchedule(user.ID, requestBody.Title, requestBody.StartTime, requestBody.EndTime)

	schedule, err := services.CreateSchedule(schedule)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	schedules, err := services.GetUserSchedules(user.ID, 1)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUserSchedules(user.ID, schedules, 1)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"schedule": schedule}
	response.SendResponse(c)
}
