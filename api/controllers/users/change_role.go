package users

import (
	"log"
	"net/http"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ChangeUserRoleRoute is used to change users role in database
func ChangeUserRoleRoute(router *gin.RouterGroup) {
	router.PUT("/changerole", ChangeRole)
}

// ChangeRole Controller for /users/changerole route
func ChangeRole(c *gin.Context) {
	val, _ := c.Get(constants.Role)
	if val != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var info serializers.ChangeRoleInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, Role are required"})
		return
	}
	info.Email = strings.ToLower(info.Email)

	status, message := serializers.EmailRegexCheck(info.Email)

	if status != http.StatusOK {
		c.JSON(status, gin.H{
			"error": message,
		})
		return
	}

	user, err := users.GetUserWithEmail(info.Email)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EmailId not in database"})
		return
	}

	serializers.ChangeRoleInfoToUserModel(&info, user)
	err = users.PatchUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Update User service error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
