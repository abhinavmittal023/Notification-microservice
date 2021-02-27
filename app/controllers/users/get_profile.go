package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetUserProfileRoute is used to get your own information
func GetUserProfileRoute(router *gin.RouterGroup) {
	router.GET("", GetUserProfile)
}

// GetUserProfile Controller for get profile/ route
func GetUserProfile(c *gin.Context) {

	userID, err := strconv.ParseUint(fmt.Sprintf("%v", c.MustGet("user_id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("String Conversion Error")
		return
	}
	user, err := users.GetUserWithID(uint64(userID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not in database"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Get user with id query error")
		return
	}

	var info serializers.UserInfo
	serializers.UserModelToUserInfo(&info, user)

	c.JSON(http.StatusOK, info)
}
