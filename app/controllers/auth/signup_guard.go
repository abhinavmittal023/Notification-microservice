package auth

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CheckIfFirstRoute used for signup guard
func CheckIfFirstRoute(router *gin.RouterGroup) {
	router.GET("/guard", CheckIfFirst)
	router.OPTIONS("/guard", preflight.Preflight)
}

// CheckIfFirst checks if another user exists to inform the front-end
func CheckIfFirst(c *gin.Context) {
	_, err := users.GetFirstUser(true)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
		return
	} else if err != nil {
		log.Println(err.Error())
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Unauthorized",
	})
}
