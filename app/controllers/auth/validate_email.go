package auth

import (
	"log"
	"net/http"
	"net/url"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/authservice"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"github.com/gin-gonic/gin"
)

// ValidateEmailRoute is used to sign in users
func ValidateEmailRoute(router *gin.RouterGroup) {
	router.GET("/token/:token", ValidateEmail)
}

// ValidateEmail Controller verifies the email after checking the token
func ValidateEmail(c *gin.Context) {
	tokenString := c.Param("token")

	userDetails, err := authservice.ValidateToken(tokenString, "validation")
	if err == authservice.ErrInvalidToken {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	userDetails.Verified = true
	err = users.PatchUser(userDetails)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	location := url.URL{Path: "http://localhost:4200/users/login"}
	c.Redirect(http.StatusFound, location.RequestURI())
}
