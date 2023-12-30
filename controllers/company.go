package controllers

import (
	"attendit/backend/models"
	"attendit/backend/services"
	"attendit/backend/services/redis"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// GetCompanyMembers godoc
// @Summary      GetCompanyMembers
// @Description  gets the company members
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        page  path  int  true  "Page"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /company/members/{page} [get]
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

// ModifyCompany godoc
// @Summary      ModifyCompany
// @Description  modifies the company
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        req  body  models.Company
// @Success      200  {object}  models.
// @Failure      400  {object}  models.
// @Router       /company [put]
func ModifyCompany(c *gin.Context) {
	company, _ := services.GetCompany()
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

// GetCompanyAttendances godoc
// @Summary      GetCompanyAttendances
// @Description  gets the company attendances
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        page  path  int  true  "Page"
// @Success      200  {object}  models.
// @Failure      400  {object}  models.
// @Router       /company/attendances/{page} [get]
func GetCompanyAttendances(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	page, _ := strconv.Atoi(c.Param("page"))
	if page == 0 {
		page = 1
	}

	attendances, err := redisServices.GetCompanyAttendancesFromCache(page)
	if err == nil {
		models.SendResponseData(c, gin.H{"attendances": attendances, "cache": true})
		return
	}

	attendances, err = services.GetAttendancesByCompany(page)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	redisServices.CacheCompanyAttendances(page, attendances)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"attendances": attendances}
	response.SendResponse(c)
}
