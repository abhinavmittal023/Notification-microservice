package users

import (
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	miscquery "code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/misc_query"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetUserRoute is used to get users from database
func GetUserRoute(router *gin.RouterGroup) {
	router.GET("/:id", GetUser)
}

// GetUser Controller for /users/:id route
func GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidID,
		})
		return
	}

	var iteratorInfo miscquery.Iterator
	err = c.BindQuery(&iteratorInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Non boolean next or prev provided",
		})
		return
	}

	if iteratorInfo.Next && iteratorInfo.Previous {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Provide either next or previous",
		})
		return
	}

	var user *models.User

	if iteratorInfo.Next {
		user, err = users.GetNextUserfromID(uint64(userID))
	} else if iteratorInfo.Previous {
		user, err = users.GetPreviousUserfromID(uint64(userID))
	} else {
		user, err = users.GetUserWithID(uint64(userID))
	}

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().IDNotInRecords})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("GetUserWithID service error", err.Error())
		return
	}

	firstRecord, err := users.GetFirstUser(false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	prevAvail := firstRecord.ID != user.ID

	lastRecord, err := users.GetLastUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	nextAvail := lastRecord.ID != user.ID

	var info serializers.UserInfo
	serializers.UserModelToUserInfo(&info, user)

	c.JSON(http.StatusOK, gin.H{
		"user_details": info,
		"next":         nextAvail,
		"previous":     prevAvail,
	})
}
