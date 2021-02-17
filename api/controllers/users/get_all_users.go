package users

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"github.com/gin-gonic/gin"
)

// GetAllUsersRoute is used to get all users from database
func GetAllUsersRoute(router *gin.RouterGroup) {
	router.GET("/", GetAllUsers)
}

// GetAllUsers Controller for get /users/ route
func GetAllUsers(c *gin.Context) {
	usersArray, err := users.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("find all users query error")
		return
	}

	var infoArray []serializers.UserInfo
	var info serializers.UserInfo

	for _, user := range usersArray {
		serializers.UserModelToUserInfo(&info, &user)
		infoArray = append(infoArray, info)
	}
	c.JSON(http.StatusOK, infoArray)
}
