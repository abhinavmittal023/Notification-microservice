package notifications

import (
	"net/smtp"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"github.com/pkg/errors"
)

// Notifications interface is used to send different types of notificaitons
type Notifications interface {
	SendNotification() error
}

// Email struct implements Notifications interface
type Email struct {
	to      string
	subject string
	message string
}

// SendNotification method send email notifications
func (email *Email) SendNotification() error {
	from := configuration.GetResp().EmailNotification.Email
	smtpHost := configuration.GetResp().EmailNotification.SMTPHost
	smtpPort := configuration.GetResp().EmailNotification.SMTPPort
	addr := smtpHost + ":" + smtpPort
	msg := []byte("Subject: " + email.subject + "\r\n" +
		"\r\n" + email.message + "\r\n")

	//  Sending email.
	err := smtp.SendMail(addr, nil, from, []string{email.to}, msg)
	if err != nil {
		return errors.Wrap(err, "Unable to send email")
	}
	return nil
}
