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

	// match, err := regexp.MatchString(constants.GetConstants().Regex.Email, info.Email)

	// if err != nil{
	// 	c.JSON(http.StatusInternalServerError, gin.H{"internal_error":"Internal Server Error"})
	// 	log.Println("Internal Server Error due to email regex")
	// 	return
	// }
	// if !match{
	// 	c.JSON(http.StatusBadRequest, gin.H{"email_invalid":"Email is invalid"})
	// 	return
	// }	
	
	// match, err = regexp.MatchString(constants.GetConstants().Regex.Password, info.Password)

	// if err != nil{
	// 	c.JSON(http.StatusInternalServerError, gin.H{"internal_error":"Internal Server Error"})
	// 	log.Println(2)
	// 	return
	// }
	// if !match{
	// 	c.JSON(http.StatusBadRequest, gin.H{"password_weak":"Password is not strong enough"})
	// 	return
	// }
	var user models.User
	serializers.LoginInfoToUserModel(info,&user)
	err := users.CreateUser(&user)
	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_error":"Internal Server Error"})
		return
	}
	to := []string{
		info.Email,
	}
	err = auth.SendValidationEmail(to,uint64(user.ID))
	if err!= nil{
		c.JSON(http.StatusBadRequest, gin.H{"smtp_error":"Email couldn't be sent"})
		return
	}
	c.JSON(http.StatusOK,gin.H{"status":"ok"})
}
