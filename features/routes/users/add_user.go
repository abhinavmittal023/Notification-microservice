package users

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/preflight"
	"github.com/gin-gonic/gin"
)

//AddUser route is used to add users to database
func AddUser(router *gin.RouterGroup) {
	router.POST("/add", users.AddUser)
	router.OPTIONS("/add", preflight.Preflight)
}
