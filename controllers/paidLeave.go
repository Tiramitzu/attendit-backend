package controllers

import (
	"attendit/backend/models"
	"attendit/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
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

	userIdHex, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))

	user, err := services.GetUserById(userId)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}
	activePaidLeave, err := services.GetActiveRequest(user.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}
	if activePaidLeave != nil {
		response.Message = "Anda masih memiliki permintaan cuti yang belum selesai"
		response.SendErrorResponse(c)
		return
	}

	paidLeave, err := services.CreatePaidLeave(user.ID, requestBody.Reason, requestBody.StartDate, requestBody.Days)
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

	userIdHex, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))

	user, err := services.GetUserById(userId)
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

// GetPaidLeaves godoc
// @Summary      GetPaidLeaves
// @Description  get all paid leave requests
// @Tags         paidLeave
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /admin/paidLeaves [get]
func GetPaidLeaves(c *gin.Context) {
	response := &models.Response{
		StatusCode: 400,
		Success:    false,
	}

	paidLeaves, err := services.GetPaidLeaves()
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	response.StatusCode = 200
	response.Success = true
	response.Data = gin.H{"paidLeaves": paidLeaves}
	response.SendResponse(c)
}

// UpdatePaidLeaveStatus godoc
// @Summary      UpdatePaidLeaveStatus
// @Description  update paid leave request status
// @Tags         paidLeave
// @Accept       json
// @Produce      json
// @Param        paidLeaveId path string true "PaidLeave ID"
// @Param        status body int true "Status"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /admin/paidLeaves/{paidLeaveId} [put]
func UpdatePaidLeaveStatus(c *gin.Context) {
	var requestBody models.PaidLeaveStatusRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)
	response := &models.Response{
		StatusCode: 400,
		Success:    false,
	}

	paidLeaveIdHex := c.Param("paidLeaveId")
	paidLeaveId, _ := primitive.ObjectIDFromHex(paidLeaveIdHex)

	status, err := strconv.Atoi(requestBody.Status)
	if err != nil {
		response.Message = "Status harus berupa angka"
		response.SendErrorResponse(c)
		return
	}

	paidLeave, err := services.UpdatePaidLeaveStatus(paidLeaveId, status)
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
