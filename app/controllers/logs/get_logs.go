package logs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
)

// GetLogsRoute is used to get all users from database
func GetLogsRoute(router *gin.RouterGroup) {
	router.GET("", GetLogs)
}

// GetLogs Controller for get /logs/ route
func GetLogs(c *gin.Context) {
	file, err := os.Open("logfile.log")

	if err != nil {
		fmt.Println("Failed to open logfile")
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()
	var logs []serializers.Logs
	var log serializers.Logs
	for _, eachline := range txtlines {
		data := []byte(eachline)
		if err := json.Unmarshal(data, &log); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		}
		logs = append(logs, log)
	}
	c.JSON(http.StatusOK, logs)
}
