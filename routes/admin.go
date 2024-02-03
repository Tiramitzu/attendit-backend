package routes

import (
	"attendit/backend/controllers"
	"attendit/backend/middlewares/validators"
	"github.com/gin-gonic/gin"
)

func AdminRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	admin := router.Group("/admin", handlers...)
	{
		admin.PUT("/user", controllers.CreateUser)
		admin.GET(
			"/user/:id",
			validators.PathIdValidator(),
			controllers.GetUser,
		)
		admin.GET("/users", controllers.GetUsers)
		admin.GET(
			"/users?page=:page",
			validators.QueryPageValidator(),
			controllers.GetUsers,
		)
		admin.PATCH(
			"/user/:id",
			validators.PathIdValidator(),
			controllers.UpdateUser,
		)
	}
}
