package middlewares

import (
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
)

// CheckIfSystemAdmin middleware checks if the user is system admin
func CheckIfSystemAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(constants.Role)
		if role != constants.SystemAdminRole {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
