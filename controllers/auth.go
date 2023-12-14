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
// @Router       /auth/register [post]
func Register(c *gin.Context) {
	var requestBody models.RegisterRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// is email in use
	err := services.CheckUserMail(requestBody.Email)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// create user record
	user, err := services.CreateUser(requestBody.UserName, requestBody.Email, requestBody.Password, requestBody.DisplayName, requestBody.Phone)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_ = c.ShouldBindBodyWith(&user, binding.JSON)

	// generate access tokens
	accessToken, err := services.GenerateAccessTokens(user)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
		"token":   accessToken.GetResponseString(),
	})
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

	// get user by email
	user, err := services.FindUserByEmail(requestBody.Email)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	// check hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		response.Message = "Email and password don't match"
		response.SendErrorResponse(c)
		return
	}

	accessToken, err := services.GetTokenById(user.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
		"token":   accessToken.GetResponseString(),
	})
}
