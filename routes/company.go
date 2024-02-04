package routes

import (
	"attendit/backend/controllers"
	"github.com/gin-gonic/gin"
)

func CompanyRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	companies := router.Group("/company", handlers...)
	{
		companies.GET(
			"",
			controllers.GetCompany,
		)
	}
}
