package auth

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/controllers/auth"
	"github.com/gin-gonic/gin"
)

//Auth route is used to create first time system admin
func Auth(router *gin.RouterGroup){
	router.GET("/token/:token",auth.ValidateEmail)
}