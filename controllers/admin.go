package controllers

import (
	"attendit/backend/models"
	"attendit/backend/services"
	redisServices "attendit/backend/services/redis"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

// GetUser godoc
// @Summary      GetUser
// @Description  gets a user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        userId path string true "User ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{userId} [get]
func GetUser(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userId, _ := c.Get("userId")

	user, err := redisServices.GetUserFromCache(userId.(primitive.ObjectID))
	if err == nil {
		response.StatusCode = http.StatusOK
		response.Success = true
		response.Data = gin.H{"user": user, "cache": true}
		response.SendResponse(c)
		return
	}

	user, err = services.GetUserById(userId.(primitive.ObjectID))

	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"user": user}
	response.SendResponse(c)
}

// GetUsers godoc
// @Summary      GetUsers
// @Description  gets users
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        page path string true "Page"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{page} [get]
func GetUsers(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	page, _ := strconv.Atoi(c.Param("page"))
	if page == 0 {
		page = 1
	}

	users, err := redisServices.GetUsersFromCache(page)
	if err == nil {
		response.StatusCode = http.StatusOK
		response.Success = true
		response.Data = gin.H{"users": users, "cache": true}
		response.SendResponse(c)
		return
	}

	users, err = services.GetUsers(page)

	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUsers(page, users)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"users": users}
	response.SendResponse(c)
}

// CreateUser godoc
// @Summary      CreateUser
// @Description  registers a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.CreateUser true "Register Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/register [put]
func CreateUser(c *gin.Context) {
	var requestBody models.CreateUser
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	err := services.CheckUserMail(requestBody.Email)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	user, err := services.CreateUser(requestBody.Email, requestBody.Password, requestBody.FullName, requestBody.Phone)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	_ = c.ShouldBindBodyWith(&user, binding.JSON)

	// generate access tokens
	accessToken, err := services.GenerateAccessTokens(user)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"user": user, "token": accessToken.GetResponseString()}
	response.SendResponse(c)
}

// UpdateUser godoc
// @Summary      UpdateUser
// @Description  updates a user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        userId path string true "User ID"
// @Param        req  body      models.ModifyUserRequest true "Update Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{userId} [patch]
func UpdateUser(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userIdHex := c.Param("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex)

	var requestBody models.ModifyUserRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	user, err := services.GetUserById(userId)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	user.Email = requestBody.Email
	user.FullName = requestBody.FullName
	user.Phone = requestBody.Phone
	user.Password = requestBody.Password

	updatedUser, err := services.UpdateUser(user)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUser(updatedUser)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"user": updatedUser}
	response.SendResponse(c)
}
