package auth

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
)

//SendValidationEmail sends validation email to new
func SendValidationEmail(to []string, userID uint64) error {

	from := configuration.GetResp().EmailNotification.Email
	password := configuration.GetResp().EmailNotification.Password
	smtpHost := configuration.GetResp().EmailNotification.SMTPHost
	smtpPort := configuration.GetResp().EmailNotification.SMTPPort
	addr := smtpHost + ":" + smtpPort

	token, err := GenerateValidationToken(userID, configuration.GetResp().Token.ExpiryTime.ValidationToken)
	if err != nil {
		log.Println("Validation Token Generation error")
		return err
	}

	link := "http://" + configuration.GetResp().Server.Domain + ":" + configuration.GetResp().Server.Port + "/api/v1/auth/token/" + token
	log.Println(link)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, err := template.ParseFiles("./shared/auth/validation_email.html")
	if err != nil {
		log.Println("Template File can't be opened")
		return err
	}
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Verify Email Address \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Link string
	}{
		Link: link,
	})

	// Sending email.
	err = smtp.SendMail(addr, auth, from, to, body.Bytes())
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent!")
	return nil
}
