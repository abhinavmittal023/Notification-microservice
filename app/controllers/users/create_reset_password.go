package users

import (
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ResetPasswordRoute is used to send reset password link to another user
func ResetPasswordRoute(router *gin.RouterGroup) {
	router.PUT("/:id/password", ResetPassword)
}

// CreatePasswordRoute is used to create a new password after checking the token
func CreatePasswordRoute(router *gin.RouterGroup) {
	router.PUT("/password/create", CreatePassword)
}

// ResetPassword is a contoller for sending reset password link
func ResetPassword(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidID,
		})
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

	message := "We have received a Request to reset your password. Do so by clicking on this link:"
	subject := "Reset Password"
	err = auth.SendHTMLEmail([]string{user.Email}, user, message, subject, true)
	if err != nil {
		log.Println("SMTP Error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// CreatePassword is the controller for PUT auth/password/create route
func CreatePassword(c *gin.Context) {

	token := c.Query("token")
	tokenHash, err := hash.Message(token, configuration.GetResp().ResetTokenHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("Hashing error for token", err.Error())
		return
	}
	user, err := users.GetUserWithToken(tokenHash)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidToken})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("GetUserWithID service error", err.Error())
		return
	}
	user.ResetToken = ""
	err = users.PatchUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("PatchUser service error", err.Error())
		return
	}

	var info serializers.ChangePasswordInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().NewPasswordRequired})
		return
	}
	newPassword, err := hash.Message(info.NewPassword, configuration.GetResp().PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("Hashing error for new password", err.Error())
		return
	}
	user.Password = newPassword
	err = users.PatchUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("PatchUser service error", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
