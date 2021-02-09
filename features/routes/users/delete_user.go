package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/preflight"
	"github.com/gin-gonic/gin"
)

//DeleteUser route is used to delete users from database
func DeleteUser(router *gin.RouterGroup) {
	router.DELETE("/delete/:id", users.DeleteUser)
	router.OPTIONS("/delete/:id", preflight.Preflight)
}
