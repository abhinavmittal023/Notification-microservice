package signin

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/signin"
	"github.com/gin-gonic/gin"
)

//SignIn route is used to sign in users
func SignIn(router *gin.RouterGroup) {
	router.POST("/", signin.SignIn)
}
