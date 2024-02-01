package routes

import (
	"attendit/backend/controllers"
	"attendit/backend/middlewares/validators"
	"github.com/gin-gonic/gin"
)

func AdminRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	admin := router.Group("/admin", handlers...)
	{
		admin.GET(
			"/user/:id",
			validators.PathIdValidator(),
			controllers.GetUser,
		)
		admin.PATCH(
			"/user/:id",
			validators.PathIdValidator(),
			controllers.VerifyUser,
		)
		admin.GET("/users", controllers.GetUsers)
		admin.GET(
			"/users/:page",
			validators.PathPageValidator(),
			controllers.GetUsers,
		)
	}
}
