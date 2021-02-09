package signup

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
)

//SignUp Controller for /signup route
func SignUp(c *gin.Context){
	var info serializers.SignupInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest,gin.H{"required":"Email,Password,FirstName are required"})
		return
	}
	info.Email = strings.ToLower(info.Email)
	info.Role = 2	//signup user will always be system admin

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
	
	info.Password = hash.Message(info.Password,configuration.GetResp().PasswordHash)

	var user models.User

	serializers.SignupInfoToUserModel(&info,&user)
	err = users.CreateUser(&user)
	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("CreateUser service error")
		return
	}
	to := []string{
		info.Email,
	}
	err = auth.SendValidationEmail(to,uint64(user.ID))
	if err!= nil{
		err = users.DeleteUserPermanently(&user)
		if err!= nil{
			c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
			log.Println("Delete User Service Error")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("SMTP Error")
		return
	}
	c.JSON(http.StatusOK,gin.H{"status":"ok"})
}
