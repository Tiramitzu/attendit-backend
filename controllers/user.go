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
// @Description  gets the current user
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{userId} [get]
func GetUser(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userIdHex, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))

	user, err := redisServices.GetUserFromCache(userId)
	if err == nil {
		response.StatusCode = http.StatusOK
		response.Success = true
		response.Data = gin.H{"user": user, "cache": true}
		response.SendResponse(c)
		return
	}

	user, err = services.GetUserById(userId)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUser(user)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"user": user}
	response.SendResponse(c)
}

// ModifyCurrentUser godoc
// @Summary      ModifyCurrentUser
// @Description  modifies the current user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        req  body      models.ModifyUserRequest    true  "Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /users/{userId} [patch]
func ModifyCurrentUser(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var requestBody models.ModifyUserRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	userIdHex, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))

	user, err := services.GetUserById(userId)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	user.Email = requestBody.Email
	user.FullName = requestBody.FullName
	user.Phone = requestBody.Phone

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

// ADMIN

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

	page, err := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}

	totalUsers, err := services.GetTotalUsers()
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	users, err := redisServices.GetUsersFromCache(page)
	if err == nil {
		response.StatusCode = http.StatusOK
		response.Success = true
		response.Data = gin.H{"users": users, "total": totalUsers, "cache": true}
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
	response.Data = gin.H{"users": users, "total": totalUsers}
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
// @Router       /admin/users [put]
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

	redisServices.CacheUser(user)

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

func DeleteUser(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userIdHex := c.Param("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdHex)

	newUsers, err := services.DeleteUser(userId)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheUsers(1, newUsers)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"users": newUsers}
	response.SendResponse(c)
}
