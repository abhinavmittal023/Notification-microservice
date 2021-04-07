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
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
)

// ChangeOwnPasswordRoute is used to change your own password
func ChangeOwnPasswordRoute(router *gin.RouterGroup) {
	router.PUT("/password", ChangePassword)
}

// ChangePassword Controller for put /users/:id/password and put /profile/password routes
func ChangePassword(c *gin.Context) {
	var userID uint64
	var err error
	var info serializers.ChangePasswordInfo
	if err := c.BindJSON(&info); err != nil {
		ve, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidDataType})
			return
		}
		var errors []string
		for _, value := range ve {
			if value.Tag() != "max" {
				errors = append(errors, fmt.Sprintf("%s is %s", value.Field(), value.Tag()))
			} else {
				errors = append(errors, fmt.Sprintf("%s cannot have more than %s characters", value.Field(), value.Param()))
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errors})
		return
	}
	if info.OldPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().OldPasswordRequired})
		return
	}
	userID, err = strconv.ParseUint(fmt.Sprintf("%v", c.MustGet(constants.ID)), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": constants.Errors().InternalError,
		})
		log.Println(err.Error())
		return
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
