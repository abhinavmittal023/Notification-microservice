package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/preflight"
	"github.com/gin-gonic/gin"
)

//GetUser route is used to get users from database
func GetUser(router *gin.RouterGroup) {
	router.GET("/get/:id", users.GetUser)
	router.OPTIONS("/get/:id", preflight.Preflight)
}
