package recipients

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/recipients"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SendRecipientCSVRoute is used to send the the recipient records in the form of csv file
func SendRecipientCSVRoute(router *gin.RouterGroup) {
	router.GET("/csv", SendRecipientCSV)
}

// SendRecipientCSV controller for get /csv route
func SendRecipientCSV(c *gin.Context) {

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

	fileBytes, err := recipients.CreateCSV(&infoArray)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=recipients.csv")
	c.Data(http.StatusOK, "text/csv", fileBytes.Bytes())

}
