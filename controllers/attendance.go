package controllers

import (
	"attendit/backend/models"
	db "attendit/backend/models/db"
	"attendit/backend/services"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AttendanceCheckIn(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var requestBody models.CheckInRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	userId, _ := c.Get("userId")

	user, err := services.FindUserById(userId.(primitive.ObjectID))
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	companyIdHex := c.Param("companyId")
	companyId, err := primitive.ObjectIDFromHex(companyIdHex)
	if err != nil {
		response.Message = "Error converting company ID"
		response.SendErrorResponse(c)
		return
	}
	company, err := services.GetCompanyById(companyId)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	getAttendance, err := services.GetAttendanceByUserAndDateAndCompany(user.ID, requestBody.Date, company.ID)

	if getAttendance != nil {
		response.Message = "You have already checked in"
		response.SendErrorResponse(c)
		return
	}

	attendance := db.NewAttendance(user.ID, company.ID, requestBody.IpAddress, requestBody.Date, requestBody.Status, requestBody.CheckIn, "")
	newAttendance, err := services.AttendanceCheckIn(attendance)

	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	c.JSON(http.StatusOK, newAttendance)
}

func AttendanceCheckOut(c *gin.Context) {
	var requestBody models.CheckOutRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	attendanceIdHex := c.Param("attendanceId")
	attendanceId, err := primitive.ObjectIDFromHex(attendanceIdHex)
	if err != nil {
		response.Message = "Error converting attendance ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	attendance, _ := services.FindAttendanceById(attendanceId)

	if attendance == nil {
		response.Message = "Attendance not found"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	attendance.CheckOut = requestBody.CheckOut
	updatedAttendance, err := services.AttendanceCheckOut(attendance)

	if err != nil {
		response.Message = "Attendance failed"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, updatedAttendance)
}

func UserAttendancesByCompany(c *gin.Context) {
	userId, _ := c.Get("userId")
	user, _ := services.FindUserById(userId.(primitive.ObjectID))
	companyId := c.Param("companyId")
	objectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}
	attendances, _ := services.FindUserAttendanceByCompany(objectId, user.ID)

	c.JSON(http.StatusOK, attendances)
}
