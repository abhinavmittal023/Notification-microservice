package auth

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/authservice"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"github.com/gin-gonic/gin"
)

// RefreshAccessTokenRoute is used to sign in users
func RefreshAccessTokenRoute(router *gin.RouterGroup) {
	router.POST("/token", RefreshAccessToken)
}

// RefreshAccessToken Provides a new access token given a valid refresh token
func RefreshAccessToken(c *gin.Context) {
	var refreshToken serializers.RefreshToken
	err := c.BindJSON(&refreshToken)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Refresh Token is Required",
		})
		return
	}

	userDetails, err := authservice.ValidateToken(refreshToken.RefreshToken, "refresh")
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

	refreshToken.AccessToken, err = auth.GenerateAccessToken(uint64(userDetails.ID), userDetails.Role, 3)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	refreshToken.RefreshToken = ""
	c.JSON(http.StatusOK, refreshToken)
}
