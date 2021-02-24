package middlewares

import (
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/notifications"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
)

// APIKeyAuth middleware checks if API Key is correct
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next() // Preflight Request
			return
		}
		authHeader := c.GetHeader(constants.Authorization)
		apiKey, err := notifications.GetAPIHash()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		if !hash.Validate(authHeader, apiKey, configuration.GetResp().APIHash) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
