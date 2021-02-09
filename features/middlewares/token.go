package middlewares

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"github.com/gin-gonic/gin"
)

//AuthorizeJWT validates and authorizes the requests
func AuthorizeJWT() gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		headerPrefix := configuration.GetResp().Token.HeaderPrefix

		if len(authHeader) > (len(headerPrefix)+2) || authHeader[:len(headerPrefix)] != headerPrefix {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		tokenString := authHeader[len(headerPrefix)+1:]

		token, err := auth.ValidateToken(tokenString)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		claims := token.Claims.(*auth.CustomClaims)
		if token.Valid && claims.TokenType == "access" {
			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}
