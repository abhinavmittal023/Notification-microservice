package users

import (
	"encoding/json"
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/features/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/services/users"
	"github.com/gin-gonic/gin"
)

//GetAllUsers Controller for /users/get route
func GetAllUsers(c *gin.Context){
	val,_ := c.Get("role")
	if val != 2{
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	usersArray,err := users.GetAllUsers()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("JSON marshalling error")
		return
	}
	c.Data(http.StatusOK, "application/json", js)
}