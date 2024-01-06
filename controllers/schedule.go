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
)

func GetUserSchedules(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userId, _ := c.Get("userId")
	user, _ := services.GetUserById(userId.(primitive.ObjectID))

	schedules, err := redisServices.GetUserSchedulesFromCache(user.ID)
	if err == nil {
		response.StatusCode = http.StatusOK
		response.Success = true
		response.Data = gin.H{"schedules": schedules, "cache": true}
		response.SendResponse(c)
		return
	}

	schedules, err = services.GetUserSchedules(user.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUserSchedules(user.ID, schedules)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"schedules": schedules}
	response.SendResponse(c)
}

func CreateUserSchedule(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var requestBody models.ScheduleRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	userId, _ := c.Get("userId")
	user, _ := services.GetUserById(userId.(primitive.ObjectID))

	schedule := &db.Schedule{
		UserId:    user.ID,
		Title:     requestBody.Title,
		StartTime: requestBody.StartTime,
		EndTime:   requestBody.EndTime,
	}

	schedule, err := services.CreateSchedule(schedule)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	schedules, err := services.GetUserSchedules(user.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUserSchedules(user.ID, schedules)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"schedule": schedule}
	response.SendResponse(c)
}
