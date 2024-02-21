package controllers

import (
	"attendit/backend/models"
	db "attendit/backend/models/db"
	"attendit/backend/services"
	redisServices "attendit/backend/services/redis"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetFeedbacks(c *gin.Context) {
	response := &models.Response{
		StatusCode: 400,
		Success:    false,
	}

	userIdHex, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))
	user, err := services.GetUserById(userId)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	var isAdmin bool
	if user.AccessLevel == 1 {
		isAdmin = true
	} else {
		isAdmin = false
	}

	feedbacks, err := redisServices.GetFeedbacksFromCache(user.ID, isAdmin)
	if err == nil {
		response.StatusCode = 200
		response.Success = true
		response.Data = gin.H{
			"feedbacks": feedbacks,
		}
		response.SendResponse(c)
		return
	}

	feedbacks, err = services.GetFeedbacks(user.ID, isAdmin)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheFeedbacks(user.ID, feedbacks, isAdmin)
	response.StatusCode = 200
	response.Success = true
	response.Data = gin.H{
		"feedbacks": feedbacks,
	}
	response.SendResponse(c)
}

func SendFeedback(c *gin.Context) {
	response := &models.Response{
		StatusCode: 400,
		Success:    false,
	}

	var requestBody models.FeedbackRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	userIdHex, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))
	user, err := services.GetUserById(userId)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	feedback := db.NewFeedback(user.ID, requestBody.Content)
	newFeedback, err := services.SendFeedback(feedback)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	feedbacksAll, err := services.GetFeedbacks(user.ID, true)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}
	redisServices.CacheFeedbacks(user.ID, feedbacksAll, true)
	feedbacksUser, err := services.GetFeedbacks(user.ID, false)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}
	redisServices.CacheFeedbacks(user.ID, feedbacksUser, false)

	response.StatusCode = 200
	response.Success = true
	response.Data = gin.H{
		"feedback": newFeedback,
	}
	response.SendResponse(c)
}
