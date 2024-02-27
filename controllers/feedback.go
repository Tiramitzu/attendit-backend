package controllers

import (
	"attendit/backend/models"
	db "attendit/backend/models/db"
	"attendit/backend/services"
	redisServices "attendit/backend/services/redis"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
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

	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}

	var isAdmin bool
	if user.AccessLevel == 1 {
		isAdmin = true
	} else {
		isAdmin = false
	}

	totalFeedbacks, err := services.GetTotalFeedbacks()
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	feedbacks, err := redisServices.GetFeedbacksFromCache(user.ID, isAdmin, page)
	if err == nil {
		response.StatusCode = 200
		response.Success = true
		response.Data = gin.H{
			"feedbacks":      feedbacks,
			"totalFeedbacks": totalFeedbacks,
		}
		response.SendResponse(c)
		return
	}

	feedbacks, err = services.GetFeedbacks(user.ID, isAdmin, page)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheFeedbacks(user.ID, feedbacks, isAdmin, page)
	response.StatusCode = 200
	response.Success = true
	response.Data = gin.H{
		"feedbacks":      feedbacks,
		"totalFeedbacks": totalFeedbacks,
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

	feedbacksAll, err := services.GetFeedbacks(user.ID, true, 1)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}
	redisServices.CacheFeedbacks(user.ID, feedbacksAll, true, 1)
	feedbacksUser, err := services.GetFeedbacks(user.ID, false, 1)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}
	redisServices.CacheFeedbacks(user.ID, feedbacksUser, false, 1)

	response.StatusCode = 200
	response.Success = true
	response.Data = gin.H{
		"feedback": newFeedback,
	}
	response.SendResponse(c)
}

func UpdateFeedbackStatus(c *gin.Context) {
	response := &models.Response{
		StatusCode: 400,
		Success:    false,
	}

	var requestBody models.FeedbackStatusRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	feedbackIdHex := c.Param("id")
	feedbackId, _ := primitive.ObjectIDFromHex(feedbackIdHex)

	feedback, err := services.UpdateFeedbackStatus(feedbackId, requestBody.Status)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	response.StatusCode = 200
	response.Success = true
	response.Data = gin.H{
		"feedback": feedback,
	}
	response.SendResponse(c)
}
