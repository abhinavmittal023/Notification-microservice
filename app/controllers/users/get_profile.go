package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetUserProfileRoute is used to get your own information
func GetUserProfileRoute(router *gin.RouterGroup) {
	router.GET("", GetUserProfile)
}

// GetUserProfile Controller for get profile/ route
func GetUserProfile(c *gin.Context) {

	userID, err := strconv.ParseUint(fmt.Sprintf("%v", c.MustGet(constants.ID)), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("String Conversion Error", err.Error())
		return
	}
	user, err := users.GetUserWithID(uint64(userID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().IDNotInRecords})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("Get user with id query error", err.Error())
		return
	}

	var info serializers.UserInfo
	serializers.UserModelToUserInfo(&info, user)

	c.JSON(http.StatusOK, info)
}
