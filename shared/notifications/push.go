package notifications

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"github.com/NaySoftware/go-fcm"
	"github.com/pkg/errors"
)

// Push struct implements Notifications interface
type Push struct {
	to    string
	title string
	body  string
}

// SendNotification method send push notifications
func (push *Push) SendNotification() error {
	var NP fcm.NotificationPayload
	NP.Title = push.title
	NP.Body = push.body

	data := map[string]string{}

	c := fcm.NewFcmClient(configuration.GetResp().PushNotification.ServerKey)
	c.NewFcmRegIdsMsg([]string{push.to}, data)
	c.SetNotificationPayload(&NP)
	status, err := c.Send()
	if status.Success != 1 || err != nil {
		return errors.Wrap(err, "Send Push Notification Error")
	}
	return nil
}
