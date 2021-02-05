package signup

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/signup"
	"github.com/gin-gonic/gin"
)

//CreateUser route is used to create first time system admin
func CreateUser(router *gin.RouterGroup){
	router.POST("/signup",signup.CreateUser)
}