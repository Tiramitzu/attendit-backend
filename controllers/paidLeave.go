package controllers

import (
	"attendit/backend/models"
	"attendit/backend/services"
	redisServices "attendit/backend/services/redis"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strconv"
	"strings"
	"time"
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

	startTime, err := time.Parse("02-01-2006", requestBody.StartDate)
	if err != nil {
		response.Message = "Format tanggal tidak sesuai"
		response.SendErrorResponse(c)
		return
	}
	startDate := primitive.NewDateTimeFromTime(startTime)
	endTime := startTime.AddDate(0, 0, requestBody.Days)
	endDate := primitive.NewDateTimeFromTime(endTime)

	if requestBody.Attachment != "" {
		// decode base64 attachment
		i := strings.Index(requestBody.Attachment, ",")
		if i < 0 {
			log.Fatal("no comma")
		}
		if !strings.Contains(requestBody.Attachment, "data:image/") {
			response.Message = "Attachment harus berupa gambar"
			response.SendErrorResponse(c)
			return
		}

		dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(requestBody.Attachment[i+1:]))
		_, err := dec.Read([]byte{})
		if err != nil {
			response.Message = "Attachment tidak valid"
			response.SendErrorResponse(c)
			return
		}
	}

	paidLeave, err := services.CreatePaidLeave(user.ID, requestBody.Reason, startDate, requestBody.Days, endDate)
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
// @Router       /users/{userId}/paidLeaves [get]
func GetPaidLeaves(c *gin.Context) {
	response := &models.Response{
		StatusCode: 400,
		Success:    false,
	}

	userIdHex, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))

	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}

	paidLeaves, err := redisServices.GetUserPaidLeavesFromCache(userId, page)
	if err == nil {
		response.StatusCode = 200
		response.Success = true
		response.Data = gin.H{"paidLeaves": paidLeaves, "cache": true}
		response.SendResponse(c)
		return
	}

	paidLeaves, err = services.GetPaidLeavesByUserId(userId, page)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUserPaidLeaves(userId, paidLeaves, page)

	response.StatusCode = 200
	response.Success = true
	response.Data = gin.H{"paidLeaves": paidLeaves}
	response.SendResponse(c)
}

// GetPaidLeavesAdmin godoc
// @Summary      GetPaidLeavesAdmin
// @Description  get all paid leave requests
// @Tags         paidLeave
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /admin/paidLeaves [get]
func GetPaidLeavesAdmin(c *gin.Context) {
	response := &models.Response{
		StatusCode: 400,
		Success:    false,
	}

	totalPaidLeaves, err := services.GetTotalPaidLeaves()
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}

	paidLeaves, err := redisServices.GetPaidLeavesFromCache(page)
	if err == nil {
		response.StatusCode = 200
		response.Success = true
		response.Data = gin.H{"paidLeaves": paidLeaves, "total": totalPaidLeaves, "cache": true}
		response.SendResponse(c)
		return
	}

	paidLeaves, err = services.GetPaidLeaves(page)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CachePaidLeaves(paidLeaves, page)

	response.StatusCode = 200
	response.Success = true
	response.Data = gin.H{"paidLeaves": paidLeaves, "total": totalPaidLeaves}
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

	paidLeaveIdHex := c.Param("id")
	paidLeaveId, _ := primitive.ObjectIDFromHex(paidLeaveIdHex)

	userIdHex, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))

	status, err := strconv.Atoi(requestBody.Status)
	if err != nil {
		response.Message = "Status harus berupa angka"
		response.SendErrorResponse(c)
		return
	}

	paidLeave, err := services.UpdatePaidLeaveStatus(paidLeaveId, status, userId)
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
