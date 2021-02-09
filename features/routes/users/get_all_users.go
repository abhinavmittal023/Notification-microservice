package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/preflight"
	"github.com/gin-gonic/gin"
)

//GetAllUsers route is used to get all users from database
func GetAllUsers(router *gin.RouterGroup) {
	router.GET("/get", users.GetAllUsers)
	router.OPTIONS("/get", preflight.Preflight)
}
