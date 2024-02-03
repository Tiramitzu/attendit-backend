package controllers

import (
	"attendit/backend/models"
	"attendit/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePaidLeave godoc
// @Summary      CreatePaidLeave
// @Description  creates a paid leave request
// @Tags         paidLeave
// @Accept       json
// @Produce      json
// @Param        reason body string true "Reason"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{userId}/paidLeave [post]
func CreatePaidLeave(c *gin.Context) {
	var requestBody models.PaidLeaveRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)
	response := &models.Response{
		StatusCode: 400,
		Success:    false,
	}

	userId, _ := c.Get("userId")

	user, err := services.GetUserById(userId.(primitive.ObjectID))
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}
	_, err = services.GetActiveRequest(user.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	paidLeave, err := services.CreatePaidLeave(user.ID, requestBody.Reason)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	response.StatusCode = 201
	response.Success = true
	response.Data = gin.H{"paidLeave": paidLeave}
	response.SendResponse(c)
}

// GetActivePaidLeave godoc
// @Summary      GetActivePaidLeave
// @Description  get active paid leave request
// @Tags         paidLeave
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{userId}/paidLeave [get]
func GetActivePaidLeave(c *gin.Context) {
	response := &models.Response{
		StatusCode: 400,
		Success:    false,
	}

	userId, _ := c.Get("userId")

	user, err := services.GetUserById(userId.(primitive.ObjectID))
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}
	paidLeave, err := services.GetActiveRequest(user.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	response.StatusCode = 200
	response.Success = true
	response.Data = gin.H{"paidLeave": paidLeave}
	response.SendResponse(c)
}
