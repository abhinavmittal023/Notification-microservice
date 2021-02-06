package auth

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"github.com/gin-gonic/gin"
)

//RefreshAccessToken Provides a new access token given a valid refresh token
func RefreshAccessToken(c *gin.Context) {

	var refreshToken serializers.RefreshToken
	err := c.BindJSON(&refreshToken)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"refresh_token_required": "Refresh Token is Required",
		})
		return
	}

	token, err := auth.ValidateToken(refreshToken.RefreshToken)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(*auth.CustomClaims)
	var userDetails models.User

	if token.Valid && claims.TokenType == "refresh" {

		err = users.GetUserWithID(&userDetails, claims.UserID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"internal_error": "Internal Server Error"})
			return
		}

		refreshToken.AccessToken, err = auth.GenerateAccessToken(claims.UserID, userDetails.Role, 3)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"internal_error": "Internal Server Error"})
			return
		}

		refreshToken.RefreshToken = ""
		c.JSON(http.StatusOK, refreshToken)

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
