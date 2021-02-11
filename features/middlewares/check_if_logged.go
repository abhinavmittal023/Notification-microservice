package middlewares

import (
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"github.com/gin-gonic/gin"
)

//CheckIfLogged middleware checks the if user was logged already in
func CheckIfLogged() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next() //Preflight Request
			return
		}
		authHeader := c.GetHeader("Authorization")
		headerPrefix := configuration.GetResp().Token.HeaderPrefix

		if len(authHeader) > (len(headerPrefix)+2) && authHeader[:len(headerPrefix)] == headerPrefix { //If token found
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": authHeader,
			})
		}
		c.Next()
	}
}
