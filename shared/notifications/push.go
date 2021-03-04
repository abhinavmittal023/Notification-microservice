package notifications

import (
	"encoding/json"
	"log"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/NaySoftware/go-fcm"
	"github.com/pkg/errors"
)

// Push struct implements Notifications interface
type Push struct {
	To    string
	Title string
	Body  string
}

// NewNotification creates fills the values in the struct with the provided ones
func (push *Push) NewNotification(to string, title string, body string) {
	push.Title = title
	push.Body = body
	push.To = to
}

// SendNotification method send push notifications
func (push *Push) SendNotification() error {
	var NP fcm.NotificationPayload
	NP.Title = push.Title
	NP.Body = push.Body

	data := map[string]string{}

	channel, err := channels.GetChannelWithType(uint(constants.ChannelIntType()[1]))
	if err != nil {
		log.Println(err.Error())
		return errors.Wrap(err, constants.Errors().InternalError)
	}

	var config serializers.PushConfig
	err = json.Unmarshal([]byte(channel.Configuration),&config)

	var c *fcm.FcmClient

	if err != nil {
		c = fcm.NewFcmClient(configuration.GetResp().PushNotification.ServerKey)
	}else{
		c = fcm.NewFcmClient(config.ServerKey)
	}	

	c.NewFcmRegIdsMsg([]string{push.To}, data)
	c.SetNotificationPayload(&NP)
	status, err := c.Send()
	if status.Success != 1 {
		return errors.New("Couldn't deliver the notification")
	} else if err != nil {
		return errors.Wrap(err, "Send Push Notification Error")
	}
	return nil
}
