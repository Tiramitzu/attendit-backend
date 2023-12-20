package models

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response Base response
type Response struct {
	StatusCode int            `json:"-"`
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

func (response *Response) SendResponse(c *gin.Context) {
	c.AbortWithStatusJSON(response.StatusCode, response)
}

func (response *Response) SendErrorResponse(c *gin.Context) {
	c.AbortWithStatusJSON(response.StatusCode, response)
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
