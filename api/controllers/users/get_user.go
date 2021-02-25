package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//GetUserRoute is used to get users from database
func GetUserRoute(router *gin.RouterGroup) {
	router.GET("/get/:id", GetUser)
	router.OPTIONS("/get/:id", preflight.Preflight)
}

//GetUser Controller for /users/get/:id route
func GetUser(c *gin.Context) {
	val, _ := c.Get(constants.Role)
	if val != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("String Conversion Error")
		return
	}
	user, err := users.GetUserWithID(uint64(userID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not in database"})
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
