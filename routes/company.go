package routes

import (
	"attendit/backend/controllers"
	"attendit/backend/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func CompanyRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	companies := router.Group("/companies", handlers...)
	{
		companies.GET(
			"/:id",
			validators.PathIdValidator(),
			controllers.GetCompany,
		)
		companies.GET(
			"/:id/members",
			validators.PathIdValidator(),
			controllers.GetCompanyMembers,
		)
		companies.GET(
			"/:id/members/:page",
			validators.PathIdValidator(),
			validators.PathPageValidator(),
			controllers.GetCompanyMembers,
		)
		companies.GET(
			"/:id/attendances/:page",
			validators.PathIdValidator(),
			validators.PathPageValidator(),
			controllers.GetCompanyAttendances,
		)
		companies.PUT(
			"/",
			validators.CreateCompanyValidator(),
			controllers.CreateCompany,
		)
		companies.PATCH(
			"/:id",
			validators.PathIdValidator(),
			controllers.ModifyCompany,
		)
		companies.DELETE(
			"/:id",
			validators.PathIdValidator(),
			controllers.DeleteCompany,
		)
	}
}
