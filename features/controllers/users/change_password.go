package users

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ChangePassword Controller for /users/changepassword route
func ChangePassword(c *gin.Context){
	userID,_ := c.Get("user_id")
	var info serializers.ChangePasswordInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest,gin.H{"required":"Email, Role are required"})
		return
	}
	info.OldPassword = hash.Message(info.OldPassword,configuration.GetResp().PasswordHash)
	info.NewPassword = hash.Message(info.NewPassword,configuration.GetResp().PasswordHash)
	
	var user models.User
	err := users.GetUserWithID(&user,userID.(uint64))
	if err == gorm.ErrRecordNotFound{
		c.JSON(http.StatusBadRequest, gin.H{"id_not_present":"Id not in database"})
		return
	}

	if info.OldPassword != user.Password{
		c.JSON(http.StatusBadRequest, gin.H{"old_password_mismatch":"Old Password is incorrect"})
		return
	}

	serializers.ChangePasswordInfoToUserModel(&info,&user)
	err = users.PatchUser(&user)
	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("Update User service error")
		return
	}
	c.JSON(http.StatusOK,gin.H{"status":"ok"})
}
