package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//CheckIfLogged middleware checks the if user was logged in to restrict making of new user when one is logged in
func CheckIfLogged() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if len(authHeader) > 10 && authHeader[:len("Bearer")] == "Bearer" { //If bearer token found
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Logout Before creating or signing into other accounts",
			})
			return
		}
		c.Next()
	}
}
