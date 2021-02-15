package recipients

import (
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/recipients"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetAllRecipientRoute is used to get recipients from database
func GetAllRecipientRoute(router *gin.RouterGroup) {
	router.GET("/get", GetAllRecipient)
	router.OPTIONS("/get", preflight.Preflight)
}

// GetAllRecipient Controller for /recipient/get route
func GetAllRecipient(c *gin.Context) {

	recipientArray, err := recipients.GetAllRecipients()
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
