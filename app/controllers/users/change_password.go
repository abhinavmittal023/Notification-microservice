package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChangeDifferentUserPasswordRoute is used to change password of another user in database
func ChangeDifferentUserPasswordRoute(router *gin.RouterGroup) {
	router.PUT("/:id/password", ChangePassword)
}

// ChangeOwnPasswordRoute is used to change your own password
func ChangeOwnPasswordRoute(router *gin.RouterGroup) {
	router.PUT("/password", ChangePassword)
}

// ChangePassword Controller for put /users/:id/password and put /profile/password routes
func ChangePassword(c *gin.Context) {
	var userID uint64
	var err error
	var info serializers.ChangePasswordInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().NewPasswordRequired})
		return
	}
	if c.Param("id") == "" {
		userID, err = strconv.ParseUint(fmt.Sprintf("%v", c.MustGet(constants.ID)), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": constants.Errors().InternalError,
			})
			log.Println(err.Error())
			return
		}
		if info.OldPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().OldPasswordRequired})
			return
		}
	} else {
		userID, err = strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": constants.Errors().InvalidID,
			})
			return
		}
	}

	info.NewPassword, err = hash.Message(info.NewPassword, configuration.GetResp().PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("Hashing error for new password", err.Error())
		return
	}

	user, err := users.GetUserWithID(uint64(userID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().IDNotInRecords})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("GetUserWithID service error", err.Error())
		return
	}

	if info.OldPassword != "" {
		match, err := hash.Validate(info.OldPassword, user.Password, configuration.GetResp().PasswordHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
			log.Println("Validation error for new password", err.Error())
			return
		}
		if !match {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.Errors().OldPasswordIncorrect})
			return
		}
	}

	serializers.ChangePasswordInfoToUserModel(&info, user)
	err = users.PatchUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("Update User service error", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
