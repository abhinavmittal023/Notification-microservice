package users

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/misc"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UpdateUserRoute is used to change users email in database
func UpdateUserRoute(router *gin.RouterGroup) {
	router.PUT("/:id/update", UpdateUser)
}

// UpdateUser Controller for put /users/:id/update route
func UpdateUser(c *gin.Context) {
	var info serializers.ChangeCredentialsInfo
	var err error
	if err = c.BindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().EmailRoleRequired})
		return
	}
	_, found := misc.FindInSlice(constants.RoleType(), info.Role)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidRole})
		return
	}

	info.Email = strings.ToLower(info.Email)

	info.ID, err = strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidID,
		})
		return
	}

	status, err := serializers.EmailRegexCheck(info.Email)

	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	testUser, err := users.GetUserWithEmail(info.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("GetUserWithEmail service error")
		return
	} else if testUser.ID != uint(info.ID) && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().EmailAlreadyPresent,
		})
	}

	user, err := users.GetUserWithID(uint64(info.ID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().IDNotInRecords})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("GetUserWithID service error")
		return
	}

	serializers.ChangeCredentialsInfoToUserModel(&info, user)
	err = users.PatchUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("Update User service error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
