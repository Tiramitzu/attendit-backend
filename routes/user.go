package routes

import (
	"attendit/backend/controllers"
	"attendit/backend/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	users := router.Group("/users/:userId", handlers...)
	{
		users.GET("", controllers.GetUser)
		users.GET(
			"/feedbacks",
			controllers.GetFeedbacks,
		)
		users.GET(
			"/attendances",
			controllers.GetUserAttendances,
		)
		users.GET(
			"/schedules",
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
		users.GET(
			"/paidLeaves",
			controllers.GetPaidLeaves,
		)
		users.POST(
			"/paidLeave",
			controllers.CreatePaidLeave,
		)
		users.POST(
			"/feedback",
			validators.FeedbackValidator(),
			controllers.SendFeedback,
		)
		users.POST(
			"/schedules",
			controllers.CreateUserSchedule,
		)
		users.POST(
			"/attendances",
			validators.CheckInValidator(),
			controllers.AttendanceCheckIn,
		)
		users.PATCH(
			"/attendances/:attendanceId",
			validators.PathAttendanceIdValidator(),
			controllers.AttendanceCheckOut,
		)
		users.PATCH("", controllers.ModifyCurrentUser)
	}
}
