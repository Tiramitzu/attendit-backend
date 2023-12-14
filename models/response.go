package models

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Response Base response
type Response struct {
	StatusCode int            `json:"-"`
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

func (response *Response) SendResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": strconv.Itoa(http.StatusOK) + ": OK"})
}

func (response *Response) SendErrorResponse(c *gin.Context) {
	c.JSON(response.StatusCode, gin.H{"message": response.Message})
}

func SendResponseData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func SendErrorResponse(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"success": "false",
		"message": message,
	})
}
