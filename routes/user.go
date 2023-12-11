package routes

import (
	"attendit/backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	users := router.Group("/users", handlers...)
	{
		users.GET("/:id", controllers.GetUser)
		users.GET("/@me", controllers.GetCurrentUser)
		users.GET("/@me/companies", controllers.GetCurrentUserCompanies)
		users.GET("/@me/attendances/:companyId", controllers.UserAttendancesByCompany)
		users.POST("/@me/attendances/:companyId", controllers.CreateAttendance)
		users.PATCH("/@me", controllers.ModifyCurrentUser)
	}
}
