package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/services/users"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//GetUser Controller for /users/get/:id route
func GetUser(c *gin.Context){
	val,_ := c.Get("role")
	if val != 2{
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID,err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("String Conversion Error")
		return
	}
	var user models.User
	err = users.GetUserWithID(&user,uint64(userID))
	if err == gorm.ErrRecordNotFound{
		c.JSON(http.StatusBadRequest, gin.H{"id_not_present":"Id not in database"})
		return
	}

	var info serializers.UserInfo
	serializers.UserModelToUserInfo(&info, &user)

	js, err := json.Marshal(&info)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("JSON marshalling error")
		return
	}
	c.Data(http.StatusOK, "application/json", js)
}