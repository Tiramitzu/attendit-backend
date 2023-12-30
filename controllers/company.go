package controllers

import (
	"attendit/backend/models"
	db "attendit/backend/models/db"
	"attendit/backend/services"
	"attendit/backend/services/redis"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCompany(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	company, err := redisServices.GetCompanyFromCache()
	if err == nil {
		models.SendResponseData(c, gin.H{"company": company, "cache": true})
		return
	}

	company, err = services.GetCompany()
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheCompany(company)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"company": company}
	response.SendResponse(c)
}

func GetCompanyMembers(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	page, _ := strconv.Atoi(c.Param("page"))
	if page == 0 {
		page = 1
	}

	users, err := redisServices.GetCompanyMembersFromCache(page)
	if err == nil {
		models.SendResponseData(c, gin.H{"users": users, "cache": true})
		return
	}

	users, err = services.GetCompanyMembers(page)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheCompanyMembers(users, page)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"users": users}
	response.SendResponse(c)
}

func ModifyCompany(c *gin.Context) {
	_ = c.ShouldBindJSON(&company)

	if company.Author != c.MustGet("userId") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": strconv.Itoa(http.StatusUnauthorized) + ": Unauthorized"})
		return
	}

	updateCompany, err := services.UpdateCompany(company)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, updateCompany)
}

	}

	}

		return
	}

	if err != nil {
		return
	}

}
