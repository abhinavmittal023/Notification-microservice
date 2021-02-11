package notifications

import (
	"log"
	"net/smtp"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
)

//SendEmail sends email message to a string of recipients
func SendEmail(to []string, subject string, message string) error {

	from := configuration.GetResp().EmailNotification.Email
	password := configuration.GetResp().EmailNotification.Password
	smtpHost := configuration.GetResp().EmailNotification.SMTPHost
	smtpPort := configuration.GetResp().EmailNotification.SMTPPort
	addr := smtpHost + ":" + smtpPort
	msg := []byte("Subject: " + subject + "\r\n" +
		"\r\n" + message + "\r\n")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(addr, auth, from, to, msg)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent Successfully!")
	return nil
}
