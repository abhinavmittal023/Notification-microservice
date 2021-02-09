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
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ChangeEmail Controller for /users/changeemail route
func ChangeEmail(c *gin.Context){
	val,_ := c.Get("role")
	if val != 2{
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var info serializers.ChangeEmailInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest,gin.H{"required":"Old Email, New Email are required"})
		return
	}
	info.OldEmail = strings.ToLower(info.OldEmail)
	info.NewEmail = strings.ToLower(info.NewEmail)

	match, err := regexp.MatchString(constants.GetConstants().Regex.Email, info.OldEmail)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("Internal Server Error due to email regex")
		return
	}
	if !match{
		c.JSON(http.StatusBadRequest, gin.H{"old_email_invalid":"Old Email is invalid"})
		return
	}

	match, err = regexp.MatchString(constants.GetConstants().Regex.Email, info.NewEmail)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("Internal Server Error due to email regex")
		return
	}
	if !match{
		c.JSON(http.StatusBadRequest, gin.H{"new_email_invalid":"New Email is invalid"})
		return
	}
	
	var user models.User
	err = users.GetUserWithEmail(&user,info.OldEmail)
	if err == gorm.ErrRecordNotFound{
		c.JSON(http.StatusBadRequest, gin.H{"email_not_present":"EmailId not in database"})
		return
	}

	serializers.ChangeEmailInfoToUserModel(&info,&user)
	err = users.PatchUser(&user)
	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("Update User service error")
		return
	}
	to := []string{
		info.NewEmail,
	}
	err = auth.SendValidationEmail(to,uint64(user.ID))
	if err!= nil{
		serializers.RevertChangesToUserEmailModel(&info,&user)
		err = users.PatchUser(&user)
		if err!= nil{
			c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
			log.Println("Revert changes Error")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("SMTP Error")
		return
	}
	c.JSON(http.StatusOK,gin.H{"status":"ok"})
}
