package routes

import (
	"attendit/backend/controllers"
	"attendit/backend/middlewares/validators"
	"github.com/gin-gonic/gin"
)

func InvitationRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	invitations := router.Group("/invitations", handlers...)
	{
		invitations.POST(
			"/",
			validators.CreateInvitationValidator(),
			controllers.CreateInvitation,
		)
		invitations.PUT(
			"/:id",
			validators.PathIdValidator(),
			validators.InsertMembersToCompanyValidator(),
			controllers.InsertMembersToCompany,
		)
	}
}
