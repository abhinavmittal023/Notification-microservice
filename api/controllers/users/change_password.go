package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChangeUserPasswordRoute is used to change users password in database
func ChangeUserPasswordRoute(router *gin.RouterGroup) {
	router.PUT("/changepassword/:id", ChangePassword)
	router.OPTIONS("/changepassword/:id", preflight.Preflight)
}

// ChangePassword Controller for /users/changepassword/:id route
func ChangePassword(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("String Conversion Error")
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
	var info serializers.ChangePasswordInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OldPassword, NewPassword are required"})
		return
	}
	if info.OldPassword != ""{
		info.OldPassword = hash.Message(info.OldPassword, configuration.GetResp().PasswordHash)
	}
	info.NewPassword = hash.Message(info.NewPassword, configuration.GetResp().PasswordHash)

	user, err := users.GetUserWithID(uint64(userID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not in database"})
		return
	}

	if info.OldPassword != "" && info.OldPassword != user.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old Password is incorrect"})
		return
	}

	serializers.ChangePasswordInfoToUserModel(&info, user)
	err = users.PatchUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Update User service error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
