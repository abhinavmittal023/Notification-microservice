package recipients

import (
	"fmt"
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/serializers"
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

	for _, recipientRecord := range *recipientRecords {

		if recipientRecord.Email != "" {
			er := serializers.EmailRegexCheck(recipientRecord.Email)

			if er == "internal_server_error" {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				log.Println("Internal Server Error due to email regex")
				return
			}
			if er == "bad_request" {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Email of ID %v is invalid", recipientRecord.ID),
				})
				return
			}
		}
		status, err := recipients.AddUpdateRecipientWithID(&recipientRecord)
		if err != nil {
			log.Println(err)
			if status == http.StatusInternalServerError {
				c.AbortWithStatusJSON(status, gin.H{"error": "Internal Server Error"})
			} else {
				c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
			}
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
