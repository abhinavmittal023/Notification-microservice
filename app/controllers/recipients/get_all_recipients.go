package recipients

import (
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/recipients"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetAllRecipientRoute is used to get recipients from database
func GetAllRecipientRoute(router *gin.RouterGroup) {
	router.GET("", GetAllRecipient)
	router.OPTIONS("", preflight.Preflight)
}

// GetAllRecipient Controller for get /recipient route
func GetAllRecipient(c *gin.Context) {

	var limit uint64
	var err error
	var offset uint64
	limitString := c.Query("limit")
	offsetString := c.Query("offset")

	if limitString != "" {
		limit, err = strconv.ParseUint(limitString, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "limit should be a unsigned integer",
			})
			return
		}
	} else {
		limit = 20
	}

	if offsetString != "" {
		offset, err = strconv.ParseUint(offsetString, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "offset should be a unsigned integer",
			})
			return
		}
	} else {
		offset = 0
	}

	recipientArray, err := recipients.GetAllRecipients(limit, offset)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	var info serializers.RecipientInfo
	var infoArray []serializers.RecipientInfo

	for _, recipient := range recipientArray {
		serializers.RecipientModelToRecipientInfo(&info, &recipient)
		infoArray = append(infoArray, info)
	}

	c.JSON(http.StatusOK, infoArray)
}
