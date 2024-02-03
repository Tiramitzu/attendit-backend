package routes

import (
	"fmt"
	"net/http"
	"time"

	"attendit/backend/docs"
	"attendit/backend/middlewares"
	"attendit/backend/models"
	"attendit/backend/services"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() *gin.Engine {
	r := gin.New()
	initRoute(r)

	r.Use(gin.LoggerWithConfig(
		gin.LoggerConfig{
			Formatter: func(param gin.LogFormatterParams) string {
				param.ClientIP = getUserIP(param.Request)
				param.TimeStamp = time.Now().In(time.FixedZone("UTC", 7*60*60))
				return fmt.Sprintf(
					"[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n",
					param.TimeStamp.Format("2006/01/02 - 15:04:05"),
					param.StatusCodeColor(),
					param.StatusCode,
					param.ResetColor(),
					param.Latency,
					param.ClientIP,
					param.MethodColor(),
					param.Method,
					param.ResetColor(),
					param.Path,
				)
			},
			Output: middlewares.LogWriter(),
		}))
	r.Use(gin.CustomRecovery(middlewares.AppRecovery()))
	r.Use(middlewares.CORSMiddleware())

	v1 := r.Group("/v1")
	{
		PingRoute(v1)
		AuthRoute(v1)
		UserRoute(v1,
			validators.PathUserIdValidator(),
			middlewares.JWTMiddleware(),
		)
		CompanyRoute(v1, middlewares.JWTMiddleware())
		AdminRoute(
			v1,
			middlewares.JWTMiddleware(),
			middlewares.IsAdminMiddleware(),
		)
	}

	docs.SwaggerInfo.BasePath = v1.BasePath() // adds /v1 to swagger base path

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func initRoute(r *gin.Engine) {
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true

	r.NoRoute(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusNotFound, c.Request.RequestURI+" not found")
	})

	r.NoMethod(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusMethodNotAllowed, c.Request.Method+" is not allowed here")
	})
}

func InitGin() {
	gin.DisableConsoleColor()
	gin.SetMode(services.Config.Mode)
	// do some other initialization staff
}

func getUserIP(httpServer *http.Request) string {
	var userIP string
	if len(httpServer.Header.Get("CF-Connecting-IP")) > 1 {
		userIP = httpServer.Header.Get("CF-Connecting-IP")
	} else if len(httpServer.Header.Get("X-Forwarded-For")) > 1 {
		userIP = httpServer.Header.Get("X-Forwarded-For")
	} else if len(httpServer.Header.Get("X-Real-IP")) > 1 {
		userIP = httpServer.Header.Get("X-Real-IP")
	} else {
		userIP = httpServer.RemoteAddr
	}

	return userIP
}
