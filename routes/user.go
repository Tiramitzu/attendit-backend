package routes

import (
	"attendit/backend/controllers"
	"attendit/backend/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	users := router.Group("/users", handlers...)
	{
		users.GET("/@me", controllers.GetCurrentUser)
		users.GET("/@me/companies", controllers.GetCurrentUserCompanies)
		users.GET(
			"/@me/attendances/:companyId",
			validators.PathCompanyIdValidator(),
			controllers.GetUserAttendancesByCompany,
		)
		users.POST(
			"/@me/attendances/:companyId",
			validators.PathCompanyIdValidator(),
			validators.CheckInValidator(),
			controllers.AttendanceCheckIn,
		)
		users.PATCH(
			"/@me/attendances/:attendanceId",
			validators.PathAttendanceIdValidator(),
			validators.CheckOutValidator(),
			controllers.AttendanceCheckOut,
		)
		users.PATCH("/@me", controllers.ModifyCurrentUser)
		users.DELETE(
			"/@me/companies/:companyId",
			validators.PathCompanyIdValidator(),
			controllers.RemoveUserFromCompany,
		)
	}
}
