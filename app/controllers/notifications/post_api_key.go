package notifications

import (
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/notifications"
	"github.com/gin-gonic/gin"
)

// PostAPIKeyRoute is used to make new API Key
func PostAPIKeyRoute(router *gin.RouterGroup) {
	router.POST("", PostAPIKey)
}

// PostAPIKey function is a controller for post notifications/api_key route
func PostAPIKey(c *gin.Context) {
	apiKey, err := notifications.PostAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"api_key": apiKey,
	})
}
