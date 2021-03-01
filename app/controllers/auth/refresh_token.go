package auth

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/authservice"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().RefreshTokenRequired,
		})
		return
	}

	userDetails, err := authservice.ValidateToken(refreshToken.RefreshToken, constants.TokenType().Refresh)
	if err == authservice.ErrInvalidToken {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	refreshToken.AccessToken, err = auth.GenerateAccessToken(uint64(userDetails.ID), userDetails.Role, configuration.GetResp().Token.ExpiryTime.AccessToken)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}
	refreshToken.RefreshToken = ""
	c.JSON(http.StatusOK, refreshToken)
}
