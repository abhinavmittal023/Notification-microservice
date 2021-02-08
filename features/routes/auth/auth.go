package auth

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/preflight"
	"github.com/gin-gonic/gin"
)

//Auth route is used to create first time system admin
func Auth(router *gin.RouterGroup) {
	router.GET("/token/:token", auth.ValidateEmail)
	router.POST("/token", auth.RefreshAccessToken)
	router.OPTIONS("/token/:token", preflight.Preflight)
	router.OPTIONS("/token", preflight.Preflight)
}
