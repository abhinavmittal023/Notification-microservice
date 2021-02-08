package signin

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/signin"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/preflight"
	"github.com/gin-gonic/gin"
)

//SignIn route is used to sign in users
func SignIn(router *gin.RouterGroup) {
	router.POST("/", signin.SignIn)
	router.OPTIONS("/", preflight.Preflight)
}
