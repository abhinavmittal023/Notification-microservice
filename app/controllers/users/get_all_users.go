package users

import (
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
)

// GetAllUsersRoute is used to get all users from database
func GetAllUsersRoute(router *gin.RouterGroup) {
	router.GET("", GetAllUsers)
}

// GetAllUsers Controller for get /users/ route
func GetAllUsers(c *gin.Context) {

	var limit uint64
	var err error
	var offset uint64
	limitString := c.Query("limit")
	offsetString := c.Query("offset")

	if limitString != "" {
		limit, err = strconv.ParseUint(limitString, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "limit should be a unsigned integer",
			})
			return
		}
	} else {
		limit = constants.DefaultLimit
	}

	if offsetString != "" {
		offset, err = strconv.ParseUint(offsetString, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "offset should be a unsigned integer",
			})
			return
		}
	} else {
		offset = constants.DefaultOffset
	}

	usersArray, err := users.GetAllUsers(limit, offset)
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
