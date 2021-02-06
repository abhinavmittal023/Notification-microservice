package signup

import (
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"github.com/gin-gonic/gin"
)

//SignUp Controller for /signup route
func SignUp(c *gin.Context){
	var info serializers.LoginInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest,gin.H{"required":"Email,Password,FirstName are required"})
		return
	}
	info.Role = 2
	var user models.User
	serializers.LoginInfoToUserModel(info,&user)
	err := users.CreateUser(&user)
	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message":"Internal Server Error"})
		return
	}
	to := []string{
		info.Email,
	}
	err = auth.SendValidationEmail(to,uint64(user.ID))
	if err!= nil{
		c.JSON(http.StatusBadRequest, gin.H{"message":"Email couldn't be sent"})
		return
	}
	c.JSON(http.StatusOK,gin.H{"status":"ok"})
}
