package controllers

import (
	"attendit/backend/models"
	"attendit/backend/services"
	"attendit/backend/services/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetCompany godoc
// @Summary      GetCompany
// @Description  gets the company
// @Tags         company
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /company [get]
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

// ModifyCompany godoc
// @Summary      ModifyCompany
// @Description  modifies the company
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        req  body  models.Company  true  "Request"
// @Success      200  {object}  models.Company
// @Failure      400  {object}  models.Company
// @Router       /admin/company [patch]
func ModifyCompany(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	company, _ := services.GetCompany()
	_ = c.ShouldBindJSON(&company)

	updateCompany, err := services.UpdateCompany(company)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheCompany(updateCompany)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"company": updateCompany}
	response.SendResponse(c)
}
