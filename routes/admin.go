package routes

import (
	"attendit/backend/controllers"
	"attendit/backend/middlewares"
	"attendit/backend/middlewares/validators"
	"github.com/gin-gonic/gin"
)

func AdminRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	admin := router.Group("/admin", handlers...)
	{
		admin.PUT("/user", controllers.CreateUser)
		admin.GET(
			"/attendances",
			validators.QueryPageValidator(),
			controllers.GetAttendances,
		)
		admin.GET(
			"/user/:id",
			validators.PathIdValidator(),
			controllers.GetUser,
		)
		admin.GET(
			"/user/:id/attendances",
			validators.PathIdValidator(),
			controllers.GetUserAttendances,
		)
		admin.GET("/users",
			validators.QueryPageValidator(),
			controllers.GetUsers)
		admin.PATCH(
			"/user/:id",
			validators.PathIdValidator(),
			controllers.UpdateUser,
		)
		admin.PATCH(
			"/company",
			middlewares.IsAdminMiddleware(),
			controllers.ModifyCompany,
		)
	}
}
