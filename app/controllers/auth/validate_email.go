package auth

import (
	"log"
	"net/http"
	"net/url"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/authservice"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ValidateEmailRoute is used to sign in users
func ValidateEmailRoute(router *gin.RouterGroup) {
	router.GET("/token/:token", ValidateEmail)
}

// ValidateEmail Controller verifies the email after checking the token
func ValidateEmail(c *gin.Context) {
	tokenString := c.Param("token")
	location := url.URL{Path: constants.LoginPath}

	userDetails, err := authservice.ValidateToken(tokenString, constants.TokenType().Validation)
	if err == authservice.ErrInvalidToken {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else if err == gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().IDNotInRecords,
		})
		return
	} else if err == authservice.ErrAlreadyVerfied {
		c.Redirect(http.StatusFound, location.RequestURI())
		return
	} else if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	userDetails.Verified = true
	err = users.PatchUser(userDetails)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}
	c.Redirect(http.StatusFound, location.RequestURI())
}
