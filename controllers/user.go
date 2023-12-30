package controllers

import (
	"net/http"
	"strconv"

	"attendit/backend/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetCurrentUser godoc
// @Summary      GetCurrentUser
// @Description  gets the current user
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/@me [get]
func GetCurrentUser(c *gin.Context) {
	userId, _ := c.Get("userId")
	user, _ := services.FindUserById(userId.(primitive.ObjectID))

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	userId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}
	user, _ := services.FindUserById(objectId)

	c.JSON(http.StatusOK, user)
}

// ModifyCurrentUser godoc
// @Summary      ModifyCurrentUser
// @Description  modifies the current user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        req  body      models.ModifyUserRequest
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/@me [patch]
func ModifyCurrentUser(c *gin.Context) {
// AttendanceCheckIn godoc
// @Summary      AttendanceCheckIn
// @Description  checks in the user
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        req  body      models.CheckInRequest
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/@me/attendances [post]
func AttendanceCheckIn(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var requestBody models.CheckInRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	company, err := services.GetCompany()
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	for _, ip := range company.IPAddresses {
		if ip == requestBody.IpAddress {
			break
		}
		response.Message = "Invalid IP Address"
		response.SendErrorResponse(c)
		return
	}

	currentDate := time.Now().Format("02-01-2006")
	currentTime := time.Now().Format("15:04:05")

	userId, _ := c.Get("userId")
	user, err := services.FindUserById(userId.(primitive.ObjectID))
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	getAttendance, err := services.GetAttendanceByUserAndDate(user.ID, currentDate)
	if getAttendance != nil {
		response.Message = "You have already checked in for today"
		response.SendErrorResponse(c)
		return
	}

	attendance := db.NewAttendance(user.ID, requestBody.IpAddress, currentDate, requestBody.Status, currentTime, "")
	newAttendance, err := services.AttendanceCheckIn(attendance)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"attendance": newAttendance}
	response.SendResponse(c)
}

// AttendanceCheckOut godoc
// @Summary      AttendanceCheckOut
// @Description  checks out the user
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        attendanceId  path  string  true  "Attendance ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/@me/attendances/{attendanceId} [patch]
func AttendanceCheckOut(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	currentTime := time.Now().Format("15:04:05")

	attendanceIdHex := c.Param("attendanceId")
	attendanceId, err := primitive.ObjectIDFromHex(attendanceIdHex)
	if err != nil {
		response.Message = "Error converting attendance ID"
		response.SendErrorResponse(c)
		return
	}

	attendance, _ := services.FindAttendanceById(attendanceId)
	if attendance == nil {
		response.Message = "Attendance not found"
		response.SendErrorResponse(c)
		return
	}
	if attendance.CheckOut != "" {
		response.Message = "You have already checked out for today"
		response.SendErrorResponse(c)
		return
	}
	attendance.CheckOut = currentTime

	updatedAttendance, err := services.AttendanceCheckOut(attendance)
	if err != nil {
		response.Message = "Attendance failed"
		response.SendErrorResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"attendance": updatedAttendance}
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
// @Router       /users/@me/attendances [get]
func GetUserAttendances(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userId, _ := c.Get("userId")
	user, _ := services.FindUserById(userId.(primitive.ObjectID))

	attendances, err := redisServices.GetUserAttendancesFromCache(user.ID)
	if err == nil {
		response.StatusCode = http.StatusOK
		response.Success = true
		response.Data = gin.H{"attendances": attendances, "cache": true}
		response.SendResponse(c)
		return
	}

	attendances, err = services.GetUserAttendances(user.ID)

	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUserAttendancesByCompany(user.ID, attendances)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"attendances": attendances}
	response.SendResponse(c)
}
