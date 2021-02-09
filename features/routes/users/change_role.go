package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/preflight"
	"github.com/gin-gonic/gin"
)

//ChangeUserRole route is used to change users role in database
func ChangeUserRole(router *gin.RouterGroup) {
	router.PUT("/changerole", users.ChangeRole)
	router.OPTIONS("/changerole", preflight.Preflight)
}
