package routes

import (
	"attendit/backend/controllers"
	"attendit/backend/middlewares/validators"
	"github.com/gin-gonic/gin"
)

func AdminRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	admin := router.Group("/admin", handlers...)
	{
		admin.PUT("/users", controllers.CreateUser)
		admin.GET("/users",
			controllers.GetUsers,
		)
		admin.GET(
			"/users/:userId",
			validators.PathUserIdValidator(),
			controllers.GetUser,
		)
		admin.GET(
			"/attendances",
			controllers.GetAttendances,
		)
		admin.GET(
			"/paidLeaves",
			controllers.GetPaidLeavesAdmin,
		)
		admin.GET(
			"/users/:userId/attendances",
			validators.PathUserIdValidator(),
			controllers.GetUserAttendances,
		)
		admin.PATCH(
			"/users/:userId",
			validators.PathUserIdValidator(),
			controllers.UpdateUser,
		)
		admin.PATCH(
			"/company",
			controllers.ModifyCompany,
		)
		admin.PATCH(
			"/paidLeaves/:id",
			validators.PathIdValidator(),
			controllers.UpdatePaidLeaveStatus,
		)
		admin.DELETE(
			"/users/:userId",
			validators.PathUserIdValidator(),
			controllers.DeleteUser,
		)
	}
}
