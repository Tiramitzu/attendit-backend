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
			controllers.GetUsers)
		admin.GET(
			"/users/:id",
			validators.PathIdValidator(),
			controllers.GetUser,
		)
		admin.GET(
			"/attendances",
			controllers.GetAttendances,
		)
		admin.GET(
			"/users/:id/attendances",
			validators.PathIdValidator(),
			controllers.GetUserAttendances,
		)
		admin.PATCH(
			"/users/:id",
			validators.PathIdValidator(),
			controllers.UpdateUser,
		)
		admin.PATCH(
			"/company",
			controllers.ModifyCompany,
		)
	}
}
