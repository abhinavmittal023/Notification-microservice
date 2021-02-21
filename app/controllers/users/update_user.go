package users

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UpdateUserRoute is used to change users email in database
func UpdateUserRoute(router *gin.RouterGroup) {
	router.PUT("/:id/update", UpdateUser)
	router.OPTIONS("/:id/update", preflight.Preflight)
}

// UpdateUser Controller for put /users/:id/update route
func UpdateUser(c *gin.Context) {
	var info serializers.ChangeCredentialsInfo
	var err error
	if err = c.BindJSON(&info); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	info.ID, err = strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "ID should be a unsigned integer",
		})
		return
	}
	if info.Email != "" {
		info.Email = strings.ToLower(info.Email)
		er := serializers.EmailRegexCheck(info.Email)

		if er == "internal_server_error" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			log.Println("Internal Server Error due to email regex")
			return
		}
		if er == "bad_request" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is invalid"})
			return
		}
	}

	user, err := users.GetUserWithID(uint64(info.ID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID not in database"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Get user with id query error")
		return
	}
	if info.Role == 0 {
		info.Role = uint(user.Role)
	} else if info.Role > constants.SystemAdminRole {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User Role provided"})
		return
	}
	if info.Email == "" {
		info.Email = user.Email
	}

	serializers.ChangeCredentialsInfoToUserModel(&info, user)
	err = users.PatchUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Update User service error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
