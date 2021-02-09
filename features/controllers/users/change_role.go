package users

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/services/users"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ChangeRole Controller for /users/changerole route
func ChangeRole(c *gin.Context){
	val,_ := c.Get("role")
	if val != 2{
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var info serializers.ChangeRoleInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest,gin.H{"required":"Email, Role are required"})
		return
	}
	info.Email = strings.ToLower(info.Email)

	match, err := regexp.MatchString(constants.GetConstants().Regex.Email, info.Email)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("Internal Server Error due to email regex")
		return
	}
	if !match{
		c.JSON(http.StatusBadRequest, gin.H{"email_invalid":"Email is invalid"})
		return
	}

	var user models.User
	err = users.GetUserWithEmail(&user,info.Email)
	if err == gorm.ErrRecordNotFound{
		c.JSON(http.StatusBadRequest, gin.H{"email_not_present":"EmailId not in database"})
		return
	}

	serializers.ChangeRoleInfoToUserModel(&info,&user)
	err = users.PatchUser(&user)
	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("Update User service error")
		return
	}
	c.JSON(http.StatusOK,gin.H{"status":"ok"})
}
