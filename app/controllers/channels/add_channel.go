package channels

import (
	"fmt"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	li "code.jtg.tools/ayush.singhal/notifications-microservice/shared/logwrapper"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/misc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
)

// AddChannelRoute is used to add channels to database
func AddChannelRoute(router *gin.RouterGroup) {
	router.POST("", AddChannel)
}

// AddChannel controller for the post channels/ route
func AddChannel(c *gin.Context) {
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
	if info.Configuration != "" {
		err := serializers.ChannelConfigValidation(&info)

		if err != nil && err.Error() == constants.Errors().InternalError {
			standardLogger.InternalServerError("Regex check in channel config validation in create channel")
			c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
			return
		} else if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	var channel models.Channel
	serializers.ChannelInfoToChannelModel(&info, &channel)

	_, err = channels.GetChannelWithName(channel.Name)
	if err != gorm.ErrRecordNotFound && err != nil {
		standardLogger.InternalServerError("Check unique channel name in create channel")
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().ChannelNamePresent})
		return
	}

	_, err = channels.GetChannelWithType(uint(info.Type))
	if err != gorm.ErrRecordNotFound && err != nil {
		standardLogger.InternalServerError("Check unique channel type in create channel")
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.Errors().ChannelTypePresent})
		return
	}

	err = channels.AddChannel(&channel)
	if err != nil {
		standardLogger.InternalServerError("Add Channel to database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.Errors().InternalError})
		return
	}

	standardLogger.EntityAdded(fmt.Sprintf("channel %s", channel.Name))
	channelInfo := serializers.ChannelModelToChannelInfo(&channel)
	c.JSON(http.StatusOK, *channelInfo)
}
