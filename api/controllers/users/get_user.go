package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetUserRoute is used to get users from database
func GetUserRoute(router *gin.RouterGroup) {
	router.GET("/get/:id", GetUser)
	router.OPTIONS("/get/:id", preflight.Preflight)
}

// GetUser Controller for /users/get/:id route
func GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "ID should be a unsigned integer",
		})
		log.Println("String Conversion Error")
		return
	}
	if userID == 0 {
		userID, err = strconv.ParseUint(fmt.Sprintf("%v", c.MustGet("user_id")), 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			log.Println("String Conversion Error")
			return
		}
	}
	user, err := users.GetUserWithID(uint64(userID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not in database"})
		return
	}

	var info serializers.UserInfo
	serializers.UserModelToUserInfo(&info, user)

	c.JSON(http.StatusOK, info)
}
