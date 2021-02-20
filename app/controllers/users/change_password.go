package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChangeDifferentUserPasswordRoute is used to change password of another user in database
func ChangeDifferentUserPasswordRoute(router *gin.RouterGroup) {
	router.PUT("/:id/password", ChangePassword)
	router.OPTIONS("/:id/password", preflight.Preflight)
}

// ChangeOwnPasswordRoute is used to change your own password
func ChangeOwnPasswordRoute(router *gin.RouterGroup) {
	router.PUT("/password", ChangePassword)
	router.OPTIONS("/password", preflight.Preflight)
}

// ChangePassword Controller for put /users/:id/password and put /profile/password routes
func ChangePassword(c *gin.Context) {
	var userID uint64
	var err error
	var info serializers.ChangePasswordInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NewPassword is required"})
		return
	}
	if c.Param("id") == "" {
		userID, err = strconv.ParseUint(fmt.Sprintf("%v", c.MustGet("user_id")), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error converting id to integer",
			})
			return
		}
		if info.OldPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OldPassword is required"})
			return
		}
	} else {
		userID, err = strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "ID should be a unsigned integer",
			})
			return
		}
	}

	info.NewPassword = hash.Message(info.NewPassword, configuration.GetResp().PasswordHash)

	user, err := users.GetUserWithID(uint64(userID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not in database"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Get user with id query error")
		return
	}

	if info.OldPassword != "" && !hash.Validate(info.OldPassword, user.Password, configuration.GetResp().PasswordHash) {
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
