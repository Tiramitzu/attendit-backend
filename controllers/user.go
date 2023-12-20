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
