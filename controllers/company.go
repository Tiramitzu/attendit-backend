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
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	companyIdHex := c.Param("id")
	companyId, err := primitive.ObjectIDFromHex(companyIdHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}

	company, err := services.GetCompanyFromCache(companyId)
	if err == nil {
		models.SendResponseData(c, gin.H{"company": company, "cache": true})
		return
	}

	company, err = services.GetCompanyById(companyId)
	if err != nil {
		response.Message = err.Error()
		response.SendErrorResponse(c)
		return
	}

	services.CacheOneCompany(company)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"company": company}
	response.SendResponse(c)
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

	invitationIdHex := c.Param("id")
	invitationId, _ := primitive.ObjectIDFromHex(invitationIdHex)
	invitation, _ := services.FindInvitationById(invitationId)

	company, _ := services.GetCompanyById(invitation.CompanyID)
	if company == nil {
		response.Message = strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	for _, member := range requestBody.Members {
		company.Members = append(company.Members, member)
	}

	_, _ = services.UpdateCompany(company)

	user, _ := services.FindUserById(invitation.UserID)
	for _, uc := range user.Companies {
		if uc == invitation.CompanyID {
			response.Message = strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}
	user.Companies = append(user.Companies, invitation.CompanyID)
	_, _ = services.UpdateUser(user)

	invitation.Status = "accepted"
	_, _ = services.UpdateInvitation(invitation)

	c.JSON(http.StatusOK, gin.H{"message": "Members inserted"})
}

func ModifyCompany(c *gin.Context) {
	companyId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}
	company, _ := services.GetCompanyById(objectId)
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
	company, _ := services.GetCompanyById(objectId)
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

func RemoveUserFromCompany(c *gin.Context) {
	userId, _ := c.Get("userId")
	user, _ := services.FindUserById(userId.(primitive.ObjectID))
	companyId := c.Param("companyId")
	objectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": strconv.Itoa(http.StatusBadRequest) + ": Invalid ID"})
		return
	}
	_ = services.RemoveUserFromCompany(objectId, user.ID)

	c.JSON(http.StatusOK, gin.H{"message": "User removed from company"})
}
