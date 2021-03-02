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
		authHeader := c.GetHeader(constants.Authorization)
		apiKey, err := notifications.GetAPIHash()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": constants.Errors().InternalError,
			})
			return
		}
		match, err := hash.Validate(authHeader, apiKey, configuration.GetResp().APIHash)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": constants.Errors().InternalError,
			})
			return
		}

		if !match {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
