package auth

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"github.com/pkg/errors"
)

// SendValidationEmail sends validation email to new
func SendValidationEmail(to []string, userID uint64) error {
	from := configuration.GetResp().EmailNotification.Email
	smtpHost := configuration.GetResp().EmailNotification.SMTPHost
	smtpPort := configuration.GetResp().EmailNotification.SMTPPort
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	token, err := GenerateValidationToken(userID, configuration.GetResp().Token.ExpiryTime.ValidationToken)
	if err != nil {
		log.Println("Validation Token Generation error")
		return errors.Wrap(err, "Unable to generate validation token")
	}

	link := fmt.Sprintf("http://%s:%s/api/v1/auth/token/%s", configuration.GetResp().Server.Domain, configuration.GetResp().Server.Port, token)

	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Unable to get working directory")
	}

	for ; strings.Split(cwd, "/")[len(strings.Split(cwd, "/"))-1] != "notifications-microservice"; cwd = filepath.Dir(cwd) {
	}

	t, err := template.ParseFiles(fmt.Sprintf("%s/shared/auth/validation_email.html", cwd))
	if err != nil {
		log.Println("Template File can't be opened")
		return errors.Wrap(err, "Unable to open template file")
	}
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Verify Email Address \n%s\n\n", mimeHeaders)))

	_ = t.Execute(&body, struct {
		Link string
	}{
		Link: link,
	})

	//  Sending email.
	err = smtp.SendMail(addr, nil, from, to, body.Bytes())
	if err != nil {
		log.Println("Unable to send email", err.Error())
		return errors.Wrap(err, "Unable to send email")
	}
	log.Println("Email Sent!")
	return nil
}
