package routes

import (
	"attendit/backend/controllers"
	"attendit/backend/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	users := router.Group("/users/:userId", handlers...)
	{
		users.GET("/", controllers.GetCurrentUser)
		users.GET(
			"/@me/attendances/:page",
			validators.PathPageValidator(),
			"/attendances",
			controllers.GetUserAttendances,
		)
		users.GET(
			"/schedules",
			controllers.GetUserSchedules,
		)
		users.GET(
			"/@me/schedules/:page",
			validators.PathPageValidator(),
			"/schedules?page=:page",
			controllers.GetUserSchedules,
		)
		users.GET(
			"/schedule/:scheduleId",
			validators.PathScheduleIdValidator(),
			controllers.GetUserSchedule,
		)
		users.GET(
			"/paidLeave",
			controllers.GetActivePaidLeave,
		)
		users.POST(
			"/@me/paidLeave",
			"/attendances",
		)
		users.POST(
			"/paidLeave",
			controllers.CreatePaidLeave,
		)
		users.POST(
			"/schedules",
			controllers.CreateUserSchedule,
		)
		users.POST(
			"/@me/attendances",
			validators.CheckInValidator(),
			controllers.AttendanceCheckIn,
		)
		users.PATCH(
			"/attendances/:attendanceId",
			validators.PathAttendanceIdValidator(),
			controllers.AttendanceCheckOut,
		)
		users.PATCH("/", controllers.ModifyCurrentUser)
	}
}
