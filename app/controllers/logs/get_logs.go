package logs

import (
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/logs"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
)

// GetLogsRoute is used to get all users from database
func GetLogsRoute(router *gin.RouterGroup) {
	router.GET("", GetLogs)
}

// GetLogs Controller for get /logs/ route
func GetLogs(c *gin.Context) {
	var err error
	var pagination serializers.Pagination
	err = c.BindQuery(&pagination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidPagination,
		})
		return
	}

	var logFilter filter.Logs
	err = c.BindQuery(&logFilter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidFilter,
		})
		return
	}

	recordsCount, err := logs.GetAllLogsCount(&logFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}
	var infoArray []serializers.Logs
	logsListResponse := serializers.LogsListResponse{
		RecordsAffected: recordsCount,
		LogInfo:         infoArray,
	}
	if recordsCount == 0 {
		c.JSON(http.StatusOK, logsListResponse)
		return
	}
	logsList, err := logs.GetAllLogs(&pagination, &logFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}
	for _, log := range logsList {
		infoArray = append(infoArray, serializers.Logs{
			Level: constants.LogLevels(log.Level),
			Time:  log.CreatedAt.String(),
			Msg:   log.Msg,
		})
	}
	logsListResponse = serializers.LogsListResponse{
		RecordsAffected: recordsCount,
		LogInfo:         infoArray,
	}

	c.JSON(http.StatusOK, logsListResponse)
}
