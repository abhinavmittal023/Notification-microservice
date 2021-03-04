package notifications

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/notifications"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
)

// GetAllNotificationsRoute is used to get all notifications from the database
func GetAllNotificationsRoute(router *gin.RouterGroup) {
	router.GET("/list", GetAllNotifications)
	router.GET("/graph", GetGraphData)
}

// GetAllNotifications function is a controller for the get notifications/ route
func GetAllNotifications(c *gin.Context) {

	var err error
	var pagination serializers.Pagination
	err = c.BindQuery(&pagination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidPagination,
		})
		return
	}

	var notificationFilter filter.Notification
	err = c.BindQuery(&notificationFilter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidFilter,
		})
		return
	}
	filter.ConvertChannelNameStringToLower(&notificationFilter)

	recordsCount, err := notifications.GetAllNotificationsCount(&notificationFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		log.Println("find all notifications count query error", err.Error())
		return
	}
	log.Println(recordsCount)
}

// GetGraphData function is a controller for get notifications/:id route
func GetGraphData(c *gin.Context) {

	var err error
	var pagination serializers.Pagination
	err = c.BindQuery(&pagination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidPagination,
		})
		return
	}

	var notificationFilter filter.Notification
	err = c.BindQuery(&notificationFilter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidFilter,
		})
		return
	}
	filter.ConvertChannelNameStringToLower(&notificationFilter)

	graphData, err := notifications.GetGraphData(&notificationFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}
	c.JSON(http.StatusOK, graphData)
}
