package users

import (
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	miscquery "code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/misc_query"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetUserRoute is used to get users from database
func GetUserRoute(router *gin.RouterGroup) {
	router.GET("/:id", GetUser)
	router.OPTIONS("/:id", preflight.Preflight)
}

// GetUser Controller for get /users/:id route
func GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "ID should be a unsigned integer",
		})
		log.Println("String Conversion Error")
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not in database"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	var info serializers.UserInfo
	serializers.UserModelToUserInfo(&info, user)

	c.JSON(http.StatusOK, info)
}
