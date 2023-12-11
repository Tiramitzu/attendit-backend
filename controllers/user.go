package controllers

import (
	"net/http"
	"strconv"

	"attendit/backend/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func ModifyCurrentUser(c *gin.Context) {
	userId, _ := c.Get("userId")
	user, _ := services.FindUserById(userId.(primitive.ObjectID))
	_ = c.ShouldBindJSON(&user)

	updateUser, err := services.UpdateUser(user)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, updateUser)
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

func CreateAttendance(c *gin.Context) {
	userId, _ := c.Get("userId")
	user, _ := services.FindUserById(userId.(primitive.ObjectID))
	companyId := c.Param("companyId")
	objectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}

	// Check if user is in the company
	companies := services.FindCompaniesByUserId(user.ID)
	for _, company := range *companies {
		if company.ID == objectId {
			attendance, _ := services.CreateAttendance(user.ID, objectId, c.ClientIP())
			c.JSON(http.StatusOK, attendance)
			return
		}
	}
}
