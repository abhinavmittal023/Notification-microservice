package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetUserRoute is used to get users from database
func GetUserRoute(router *gin.RouterGroup) {
	router.GET("/get/:id", GetUser)
}

// GetUser Controller for /users/get/:id route
func GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "ID should be a unsigned integer",
		})
		return
	}
	if userID == 0 {
		userID, err = strconv.Atoi(fmt.Sprintf("%v", c.MustGet("user_id")))
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("GetUserWithID service error")
		return
	}

	var info serializers.UserInfo
	serializers.UserModelToUserInfo(&info, user)

	js, err := json.Marshal(&info)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("JSON marshalling error")
		return
	}
	c.Data(http.StatusOK, "application/json", js)
}
