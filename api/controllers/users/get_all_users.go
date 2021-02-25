package users

import (
	"encoding/json"
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
)

// GetAllUsersRoute is used to get all users from database
func GetAllUsersRoute(router *gin.RouterGroup) {
	router.GET("/get", GetAllUsers)
}

// GetAllUsers Controller for /users/get route
func GetAllUsers(c *gin.Context) {
	val, _ := c.Get(constants.Role)
	if val != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

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

	js, err := json.Marshal(&infoArray)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("JSON marshalling error")
		return
	}
	c.Data(http.StatusOK, "application/json", js)
}
