package signup

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/signup"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/preflight"
	"github.com/gin-gonic/gin"
)

//SignUp route is used to create first time system admin
func SignUp(router *gin.RouterGroup) {
	router.POST("/", signup.SignUp)
	router.OPTIONS("/", preflight.Preflight)
}
