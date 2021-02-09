package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/preflight"
	"github.com/gin-gonic/gin"
)

//ChangeUserEmail route is used to change users email in database
func ChangeUserEmail(router *gin.RouterGroup) {
	router.PUT("/changeemail", users.ChangeEmail)
	router.OPTIONS("/changeemail", preflight.Preflight)
}
