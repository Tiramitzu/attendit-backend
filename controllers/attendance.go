package controllers

import (
	"attendit/backend/models"
	db "attendit/backend/models/db"
	"attendit/backend/services"
	redisServices "attendit/backend/services/redis"
	"fmt"
	"github.com/Kagami/go-face"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const dataDir = "testdata"

var (
	modelsDir = filepath.Join(dataDir, "models")
	imagesDir = filepath.Join(dataDir, "images")
)

// AttendanceCheckIn godoc
// @Summary      AttendanceCheckIn
// @Description  checks in the user
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        ipAddress  body  string  true  "IP Address"
// @Param        status     body  string  true  "Status"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{userId}/attendances [post]
func AttendanceCheckIn(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	isIpSame := false

	var requestBody models.CheckInRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	company, err := services.GetCompany()
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	ipAddress, err := services.GetClientIP(c.Request)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	for _, ip := range company.IPAddresses {
		if ip == ipAddress {
			isIpSame = true
			break
		}
	}

	if !isIpSame {
		response.Message = "Anda tidak diizinkan untuk melakukan absensi dari alamat IP ini."
		response.SendErrorResponse(c)
		return
	}

	userIdHex, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))

	currentDate := time.Now().Format("02-01-2006")
	currentTime := time.Now().Format("15:04:05")

	user, err := services.GetUserById(userId)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	getAttendance, err := services.GetAttendanceByUserAndDate(user.ID, currentDate)
	if getAttendance != nil {
		response.Message = "Anda sudah melakukan absensi untuk hari ini."
		response.SendErrorResponse(c)
		return
	}

	// Init the recognizer.
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	defer rec.Close()

	fmt.Println("Recognizer Initialized")

	photoPath, err := services.SaveImage(user.Photo, user, imagesDir, "photos")

	// Print success message
	log.Println("Image saved successfully at:", photoPath)
	faces, err := rec.RecognizeFile(photoPath)
	if err != nil {
		err = os.Remove(photoPath)
		if err != nil {
			response.Message = err.Error()
			response.SendErrorResponse(c)
		}

		response.Message = err.Error()
		response.SendErrorResponse(c)
	}
	if len(faces) < 1 {
		err = os.Remove(photoPath)
		if err != nil {
			response.Message = err.Error()
			response.SendErrorResponse(c)
		}

		response.Message = "Wajah tidak ditemukan"
		response.SendErrorResponse(c)
		return
	}

	var samples []face.Descriptor
	var cats []int32
	for i, f := range faces {
		samples = append(samples, f.Descriptor)
		cats = append(cats, int32(i))
	}
	// Name the categories, i.e. people on the image.
	labels := []string{
		user.FullName,
	}
	// Pass samples to the recognizer.
	rec.SetSamples(samples, cats)

	imagesPath, err := services.SaveImage(requestBody.Image, user, imagesDir, "images")

	// Recognize faces on that image.
	imageFaces, err := rec.RecognizeFile(imagesPath)
	if err != nil {
		err = os.Remove(imagesPath)
		if err != nil {
			response.Message = err.Error()
			response.SendErrorResponse(c)
		}

		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	if len(imageFaces) < 1 {
		err = os.Remove(imagesPath)
		if err != nil {
			response.Message = err.Error()
			response.SendErrorResponse(c)
		}

		response.Message = "Wajah tidak ditemukan"
		response.SendErrorResponse(c)
		return
	}

	// Classify the image
	catID := rec.Classify(imageFaces[0].Descriptor)
	if catID < 0 {
		response.Message = "Wajah tidak ditemukan"
		response.SendErrorResponse(c)
		return
	}

	fmt.Println(labels[catID])

	attendance := db.NewAttendance(user.ID, ipAddress, currentDate, requestBody.Status, currentTime, "")
	newAttendance, err := services.AttendanceCheckIn(attendance)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	err = os.Remove(photoPath)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
	}

	err = os.Remove(imagesPath)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
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
// @Router       /users/{userId}/attendances/{attendanceId} [patch]
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

	attendance, _ := services.GetAttendanceById(attendanceId)
	if attendance == nil {
		response.Message = "Absensi tidak ditemukan"
		response.SendErrorResponse(c)
		return
	}
	if attendance.CheckOut != "" {
		response.Message = "Anda telah melakukan absen keluar untuk hari ini"
		response.SendErrorResponse(c)
		return
	}
	attendance.CheckOut = currentTime

	updatedAttendance, err := services.AttendanceCheckOut(attendance)
	if err != nil {
		response.Message = "Absen keluar gagal."
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
// @Param        userId path string true "User ID"
// @Param        page query int false "Page"
// @Param        from query string false "From Date"
// @Param        to query string false "To Date"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{userId}/attendances [get]
func GetUserAttendances(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	page, _ := strconv.Atoi(c.Query("page"))
	userIdHex, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))
	if page == 0 {
		page = 1
	}

	fromDate := c.Query("from")
	toDate := c.Query("to")

	if fromDate != "" && toDate != "" {
		attendances, err := services.GetUserAttendancesByDate(userId, fromDate, toDate, page)
		if err != nil {
			response.Message = err.Error()
			response.SendErrorResponse(c)
			return
		}

		response.StatusCode = http.StatusOK
		response.Success = true
		response.Data = gin.H{"attendances": attendances}
		response.SendResponse(c)
		return
	}

	user, err := services.GetUserById(userId)
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

// GetAttendances godoc
// @Summary      GetAttendances
// @Description  gets the company attendances
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        page query string true "Page"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /admin/attendances [get]
func GetAttendances(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}

	fromDate := c.Query("from")
	toDate := c.Query("to")

	if fromDate != "" && toDate != "" {
		totalAttendances, err := redisServices.GetAttendanceTotalByDateFromCache(fromDate, toDate)
		if err != nil {
			totalAttendances, err = services.GetTotalAttendancesByDate(fromDate, toDate)
			if err != nil {
				response.Message = err.Error()
				response.SendErrorResponse(c)
				return
			}

			redisServices.CacheAttendanceTotalByDate(fromDate, toDate, totalAttendances)
		}

		attendances, err := redisServices.GetAttendancesByDateFromCache(fromDate, toDate, page)
		if err == nil {
			response.StatusCode = http.StatusOK
			response.Success = true
			response.Data = gin.H{"attendances": attendances, "totals": totalAttendances, "cache": true}
			response.SendResponse(c)
			return
		}

		attendances, err = services.GetAttendancesByDate(fromDate, toDate, page)

		redisServices.CacheAttendancesByDate(fromDate, toDate, page, attendances)

		response.StatusCode = http.StatusOK
		response.Success = true
		response.Data = gin.H{"attendances": attendances, "totals": totalAttendances}
		response.SendResponse(c)
		return
	}

	totalAttendances, err := redisServices.GetAttendanceTotalFromCache()
	if err != nil {
		totalAttendances, err = services.GetTotalAttendances()
		if err != nil {
			response.Message = err.Error()
			response.SendErrorResponse(c)
			return
		}

		redisServices.CacheAttendanceTotal(totalAttendances)
	}

	attendances, err := redisServices.GetAttendancesFromCache(page)
	if err == nil {
		models.SendResponseData(c, gin.H{"attendances": attendances, "totals": totalAttendances, "cache": true})
		return
	}

	attendances, err = services.GetAttendances(page)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheAttendances(page, attendances)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"attendances": attendances, "totals": totalAttendances}
	response.SendResponse(c)
}
