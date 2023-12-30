package routes

import (
	"attendit/backend/controllers"
	"attendit/backend/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func CompanyRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	companies := router.Group("/company", handlers...)
	{
		companies.GET(
			"",
			validators.PathIdValidator(),
			controllers.GetCompany,
		)
		companies.GET(
			"/members",
			validators.PathIdValidator(),
			controllers.GetCompanyMembers,
		)
		companies.GET(
			"/members/:page",
			validators.PathIdValidator(),
			validators.PathPageValidator(),
			controllers.GetCompanyMembers,
		)
		companies.GET(
			"/attendances",
			validators.PathIdValidator(),
			controllers.GetCompanyAttendances,
		)
		companies.GET(
			"/attendances/:page",
			validators.PathIdValidator(),
			validators.PathPageValidator(),
			controllers.GetCompanyAttendances,
		)
		companies.PATCH(
			"",
			validators.PathIdValidator(),
			controllers.ModifyCompany,
		)
	}
}
