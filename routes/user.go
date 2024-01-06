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
		users.GET(
			"/@me/attendances",
			controllers.GetUserAttendances,
		)
		users.GET(
			"/@me/attendances/:page",
			validators.PathPageValidator(),
			controllers.GetUserAttendances,
		)
		users.GET(
			"/@me/schedules",
			controllers.GetUserSchedules,
		)
		users.GET(
			"/@me/schedules/:page",
			validators.PathPageValidator(),
			controllers.GetUserSchedules,
		)
		users.GET(
			"/@me/schedules/:scheduleId",
			validators.PathScheduleIdValidator(),
			controllers.GetUserSchedule,
		)
		users.POST(
			"/@me/schedules",
			controllers.CreateUserSchedule,
		)
		users.POST(
			"/@me/attendances",
			validators.CheckInValidator(),
			controllers.AttendanceCheckIn,
		)
		users.PATCH(
			"/@me/attendances/:attendanceId",
			validators.PathAttendanceIdValidator(),
			controllers.AttendanceCheckOut,
		)
		users.PATCH("/@me", controllers.ModifyCurrentUser)
	}
}
