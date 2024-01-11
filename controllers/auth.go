package controllers

import (
	"attendit/backend/models"
	"attendit/backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary      Register
// @Description  registers a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.RegisterRequest true "Register Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/register [put]
func Register(c *gin.Context) {
	var requestBody models.RegisterRequest
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

// Login godoc
// @Summary      Login
// @Description  login a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.LoginRequest true "Login Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var requestBody models.LoginRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	user, err := services.GetUserByEmail(requestBody.Email)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		response.Message = "Email and password don't match"
		response.SendErrorResponse(c)
		return
	}

	if user.IsVerified == false {
		response.Message = "Please contact your admin to verify your account"
		response.SendErrorResponse(c)
		return
	}

	accessToken, err := services.GetTokenById(user.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"user": user, "token": accessToken.GetResponseString()}
	response.SendResponse(c)
}
