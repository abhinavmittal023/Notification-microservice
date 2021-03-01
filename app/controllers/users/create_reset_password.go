package users

import (
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ResetPasswordRoute is used to send reset password link to another user
func ResetPasswordRoute(router *gin.RouterGroup) {
	router.PUT("/:id/password", ChangePassword)
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
	err = auth.SendHTMLEmail([]string{user.Email}, user, message, true)
	if err != nil {
		log.Println("SMTP Error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
