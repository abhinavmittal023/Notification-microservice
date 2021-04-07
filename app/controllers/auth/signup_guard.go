package auth

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CheckIfFirstRoute used for signup guard
func CheckIfFirstRoute(router *gin.RouterGroup) {
	router.GET("/guard", CheckIfFirst)
}

// CheckIfFirst checks if another user exists to inform the front-end
func CheckIfFirst(c *gin.Context) {
	_, err := users.GetFirstUser(true)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{
			"status": "Allowed",
		})
		return
	} else if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": constants.Errors().InternalError,
		})
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"status": "Not Allowed",
	})
}
