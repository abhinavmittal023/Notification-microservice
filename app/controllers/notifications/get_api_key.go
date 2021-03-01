package notifications

import (
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/notifications"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
)

// GetAPIKeyRoute is used to make new API Key
func GetAPIKeyRoute(router *gin.RouterGroup) {
	router.GET("/new", GetAPIKey)
}

// GetAPIKey function is a controller for post notifications/api_key route
func GetAPIKey(c *gin.Context) {
	apiKey, err := notifications.GetAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.Errors().InternalError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"api_key": apiKey,
	})
}
