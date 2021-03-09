package notifications

import (
	"encoding/json"
	"log"
	"net/smtp"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/pkg/errors"
)

// Notifications interface is used to send different types of notificaitons
type Notifications interface {
	SendNotification() error
	NewNotification(to string, title string, body string)
}

// Email struct implements Notifications interface
type Email struct {
	To      string
	Subject string
	Message string
}

// NewNotification creates fills the values in the struct with the provided ones
func (email *Email) NewNotification(to string, title string, body string) {
	email.Message = body
	email.To = to
	email.Subject = title
}

// SendNotification method send email notifications
func (email *Email) SendNotification() error {
	channel, err := channels.GetChannelWithType(uint(constants.ChannelIntType()[0]))
	if err != nil {
		log.Println(err.Error())
		return errors.Wrap(err, constants.Errors().InternalError)
	}
	var from, fromEmail, smtpHost, smtpPort string
	var password string
	var config serializers.EmailConfig
	err = json.Unmarshal([]byte(channel.Configuration), &config)
	if err != nil || config.Email == "" || config.SMTPHost == "" || config.SMTPPort == "" {
		fromEmail = configuration.GetResp().EmailNotification.Email
		from = configuration.GetResp().EmailNotification.From
		password = configuration.GetResp().EmailNotification.Password
		smtpHost = configuration.GetResp().EmailNotification.SMTPHost
		smtpPort = configuration.GetResp().EmailNotification.SMTPPort
	} else {
		fromEmail = strings.ToLower(config.Email)
		from = config.From
		if from == ""{
			from = fromEmail
		}
		password = config.Password
		smtpHost = config.SMTPHost
		smtpPort = config.SMTPPort
	}
	addr := smtpHost + ":" + smtpPort
	
	msg := []byte("Subject: " + email.Subject + "\r\n" +
		"\r\n" + email.Message + "\r\n")

	// Authentication
	auth := smtp.PlainAuth("", fromEmail, password, addr)

	//  Sending email.
	err = smtp.SendMail(addr, auth, from, []string{email.To}, msg)
	if err != nil {
		return errors.Wrap(err, "Unable to send email")
	}
	return nil
}
