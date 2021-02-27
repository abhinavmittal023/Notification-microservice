package users

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChangeOwnPasswordRoute is used to change your own password in database
func ChangeOwnPasswordRoute(router *gin.RouterGroup) {
	router.PUT("/changepassword", ChangePassword)
}

// ChangePassword Controller for profile/changepassword route
func ChangePassword(c *gin.Context) {
	userID := c.MustGet(constants.ID)
	var info serializers.ChangePasswordInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old and New Password are required"})
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
