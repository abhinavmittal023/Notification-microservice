package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/preflight"
	"github.com/gin-gonic/gin"
)

//ChangeUserPassword route is used to change users password in database
func ChangeUserPassword(router *gin.RouterGroup) {
	router.PUT("/changepassword", users.ChangePassword)
	router.OPTIONS("/changepassword", preflight.Preflight)
}
