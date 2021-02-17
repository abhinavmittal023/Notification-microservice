package recipients

import (
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/controllers/preflight"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/recipients"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetRecipientRoute is used to get recipients from database
func GetRecipientRoute(router *gin.RouterGroup) {
	router.GET("/:id", GetRecipient)
	router.OPTIONS("/:id", preflight.Preflight)
}

// GetRecipient Controller for get /recipient/:id route
func GetRecipient(c *gin.Context) {
	recipientID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("String Conversion Error")
		return
	}
	recipient, err := recipients.GetRecipientWithID(uint64(recipientID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not in database"})
		return
	}
	var info serializers.RecipientInfo
	serializers.RecipientModelToRecipientInfo(&info, recipient)

	c.JSON(http.StatusOK, info)
}
