package auth

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"github.com/pkg/errors"
)

// SendValidationEmail sends validation email to new
func SendValidationEmail(to []string, userID uint64) error {
	from := configuration.GetResp().EmailNotification.Email
	//password := configuration.GetResp().EmailNotification.Password
	smtpHost := configuration.GetResp().EmailNotification.SMTPHost
	smtpPort := configuration.GetResp().EmailNotification.SMTPPort
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	token, err := GenerateValidationToken(userID, configuration.GetResp().Token.ExpiryTime.ValidationToken)
	if err != nil {
		log.Println("Validation Token Generation error")
		return errors.Wrap(err, "Unable to generate validation token")
	}

	link := fmt.Sprintf("http://%s:%s/api/v1/auth/token/%s", configuration.GetResp().Server.Domain, configuration.GetResp().Server.Port, token)

	//  Authentication.
	//auth := smtp.PlainAuth("", from, password, smtpHost)

	t, err := template.ParseFiles("./shared/auth/validation_email.html")
	if err != nil {
		log.Println("Template File can't be opened")
		return errors.Wrap(err, "Unable to open template file")
	}
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Verify Email Address \n%s\n\n", mimeHeaders)))

	err = t.Execute(&body, struct {
		Link string
	}{
		Link: link,
	})
	if err != nil {
		log.Println("Unable to write to template")
		return errors.Wrap(err, "Unable to write to template")
	}

	//  Sending email.
	err = smtp.SendMail(addr, nil, from, to, body.Bytes())
	if err != nil {
		log.Println("Unable to send email")
		return errors.Wrap(err, "Unable to send email")
	}
	return nil
}
