package users

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ChangeUserPasswordRoute is used to change users password in database
func ChangeUserPasswordRoute(router *gin.RouterGroup) {
	router.PUT("/changepassword", ChangePassword)
	router.OPTIONS("/changepassword", preflight.Preflight)
}

//ChangePassword Controller for /users/changepassword route
func ChangePassword(c *gin.Context) {
	userID, _ := c.Get(constants.ID)
	var info serializers.ChangePasswordInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, Role are required"})
		return
	}

	var err error
	info.OldPassword, err = hash.Message(info.OldPassword, configuration.GetResp().PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Hashing error for old password")
		return
	}
	info.NewPassword, err = hash.Message(info.NewPassword, configuration.GetResp().PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Hashing error for new password")
		return
	}

	user, err := users.GetUserWithID(userID.(uint64))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not in database"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("GetUserWithID service error")
		return
	}

	if info.OldPassword != user.Password {
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
