package routes

import (
	"attendit/backend/controllers"
	"attendit/backend/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.PUT(
			"/register",
			validators.RegisterValidator(),
			controllers.Register,
		)

		auth.POST(
			"/login",
			validators.LoginValidator(),
			controllers.Login,
		)
	}
}
