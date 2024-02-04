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
			validators.QueryPageValidator(),
			controllers.GetUsers)
		admin.GET(
			"/users/:id",
			validators.PathIdValidator(),
			controllers.GetUser,
		)
		admin.GET(
			"/attendances",
			validators.QueryPageValidator(),
			controllers.GetAttendances,
		)
		admin.GET(
			"/users/:id/attendances",
			validators.QueryPageValidator(),
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
