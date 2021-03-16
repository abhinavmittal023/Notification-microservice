package channels

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/logs"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	li "code.jtg.tools/ayush.singhal/notifications-microservice/shared/logwrapper"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/misc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
)

// UpdateChannelRoute is used to update existing channels
func UpdateChannelRoute(router *gin.RouterGroup) {
	router.PUT(":id", UpdateChannel)
}

// UpdateChannel controller for put the channels/:id route
func UpdateChannel(c *gin.Context) {
	f, err := li.OpenFile()
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	}
	defer f.Close()
	var standardLogger = li.NewLogger()
	standardLogger.SetOutput(f)
	var info serializers.ChannelInfo
	if err := c.BindJSON(&info); err != nil {
		ve, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidDataType})
			return
		}
		var errors []string
		for _, value := range ve {
			if value.Tag() != "max" {
				errors = append(errors, fmt.Sprintf("%s is %s", value.Field(), value.Tag()))
			} else {
				errors = append(errors, fmt.Sprintf("%s cannot have more than %s characters", value.Field(), value.Param()))
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errors})
		return
	}
	_, found := misc.FindInSlice(constants.ChannelIntType(), int(info.Type))
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidType})
		return
	}
	if info.Priority > constants.MaxPriority {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().InvalidPriority})
		return
	}

	channelID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().InvalidID,
		})
		return
	}

	channel, err := channels.GetChannelWithID(uint(channelID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.Errors().IDNotInRecords,
		})
		return
	} else if err != nil {
		standardLogger.InternalServerError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	testChannel, err := channels.GetChannelWithName(strings.ToLower(info.Name))
	if err != gorm.ErrRecordNotFound && err != nil {
		standardLogger.InternalServerError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	} else if testChannel.ID != channel.ID && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().ChannelNamePresent})
		return
	}

	testChannel, err = channels.GetChannelWithType(uint(info.Type))
	if err != gorm.ErrRecordNotFound && err != nil {
		standardLogger.InternalServerError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	} else if testChannel.ID != channel.ID && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().ChannelTypePresent})
		return
	}

	if info.Configuration != "" {
		err := serializers.ChannelConfigValidation(&info)

		if err != nil && err.Error() == constants.Errors().InternalError {
			standardLogger.InternalServerError(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
			return
		} else if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	serializers.ChannelInfoToChannelModel(&info, channel)

	err = channels.PatchChannel(channel)
	if err != nil {
		standardLogger.InternalServerError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	logs.AddLogs(constants.InfoLog, fmt.Sprintf("Channel %s updated", channel.Name))
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
