package controllers

import (
	"attendit/backend/models"
	db "attendit/backend/models/db"
	"attendit/backend/services"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCurrentUserCompanies(c *gin.Context) {
	userId, _ := c.Get("userId")
	companies := services.FindCompaniesByUserId(userId.(primitive.ObjectID))

	c.JSON(http.StatusOK, companies)
}

func GetCompany(c *gin.Context) {
	companyId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}
	company, _ := services.FindCompanyById(objectId)

	c.JSON(http.StatusOK, company)
}

func GetCompanyMembers(c *gin.Context) {
	companyId := c.Param("id")
	page, _ := strconv.Atoi(c.Param("page"))
	if page == 0 {
		page = 1
	}
	objectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}
	users := services.FindMembersByCompanyId(objectId, page)

	c.JSON(http.StatusOK, users)
}

func GetCompanyAttendances(c *gin.Context) {
	companyId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}
	attendances, _ := services.FindAttendanceByCompany(objectId)

	c.JSON(http.StatusOK, attendances)
}

func CreateCompany(c *gin.Context) {
	var requestBody models.CreateCompanyRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	authorId, _ := primitive.ObjectIDFromHex(requestBody.Author)

	currentUser, _ := services.FindUserById(authorId)

	if currentUser == nil {
		response.Message = strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	member := []db.Member{
		{
			Name:  currentUser.DisplayName,
			Email: currentUser.Email,
			Phone: currentUser.Phone,
			Role:  "admin",
		},
	}

	company := db.NewCompany(authorId, requestBody.Name, requestBody.IPAddresses, requestBody.CheckInTime, requestBody.CheckOutTime, member)
	newCompany, err := services.CreateCompany(company)

	currentUser.Companies = append(currentUser.Companies, newCompany.ID)
	_, _ = services.UpdateUser(currentUser)

	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, newCompany)
}

func InsertMembersToCompany(c *gin.Context) {
	var requestBody models.InsertMembersToCompanyRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	companyId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		response.Message = strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	company, _ := services.FindCompanyById(objectId)
	if company == nil {
		response.Message = strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if company.Author != c.MustGet("userId") {
		response.Message = strconv.Itoa(http.StatusUnauthorized) + ": Unauthorized"
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	err = services.InsertMembersToCompany(objectId, requestBody.Members)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Members inserted"})
}

func ModifyCompany(c *gin.Context) {
	companyId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}
	company, _ := services.FindCompanyById(objectId)
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

func DeleteCompany(c *gin.Context) {
	companyId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}
	company, _ := services.FindCompanyById(objectId)
	if company == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}

	if company.Author != c.MustGet("userId") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": strconv.Itoa(http.StatusUnauthorized) + ": Unauthorized"})
		return
	}

	err = services.DeleteCompany(objectId)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted"})
}
