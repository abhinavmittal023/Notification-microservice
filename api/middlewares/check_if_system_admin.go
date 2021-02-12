package middlewares

import (
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
)

// CheckIfSystemAdmin middleware checks if the user is system admin
func CheckIfSystemAdmin() gin.HandlerFunc{
	return func(c *gin.Context){
		if c.Request.Method == "OPTIONS" {
			c.Next() // Preflight Request
			return
		}
		role, _ := c.Get("role")
		if role != constants.SystemAdminRole {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}