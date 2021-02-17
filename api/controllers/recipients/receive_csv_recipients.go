package recipients

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/recipients"
	"github.com/gin-gonic/gin"
)

// AddUpdateRecipientRoute is used to allow creation and updation of recipients from csv
func AddUpdateRecipientRoute(router *gin.RouterGroup) {
	router.POST("/csv", AddUpdateRecipient)
}

// AddUpdateRecipient controller for post /recipient/csv route
func AddUpdateRecipient(c *gin.Context) {

	rFile, err := c.FormFile("recipients")

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "File Format error",
		})
		return
	}
	recipientRecords, err := recipients.ReadCSV(rFile)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid CSV file",
		})
		return
	}

	status, errors := recipients.AddUpdateRecipients(recipientRecords)

	if status == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	} else {
		c.AbortWithStatusJSON(status, errors)
	}
}
