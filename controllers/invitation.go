package controllers

import (
	"attendit/backend/models"
	db "attendit/backend/models/db"
	"attendit/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

func CreateInvitation(c *gin.Context) {
	var requestBody models.CreateInvitationRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userIdHex, exist := c.Get("userId")
	if !exist {
		response.Message = strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))

	company, _ := services.GetCompanyById(requestBody.CompanyID)
	if company == nil {
		response.Message = strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	for _, member := range company.Members {
		if member.Role != "admin" && member.ID == userId {
			response.Message = strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	user, _ := services.FindUserById(requestBody.UserID)
	if user == nil {
		response.Message = strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	invitation := db.NewInvitation(userId, user.ID, company.ID, requestBody.Role)
	newInvitation, err := services.CreateInvitation(invitation)

	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user.Invitations = append(user.Invitations, newInvitation.ID)
	_, _ = services.UpdateUser(user)

	c.JSON(http.StatusOK, newInvitation)
}
