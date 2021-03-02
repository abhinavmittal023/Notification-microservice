package notifications

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
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

	c := fcm.NewFcmClient(configuration.GetResp().PushNotification.ServerKey)
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
