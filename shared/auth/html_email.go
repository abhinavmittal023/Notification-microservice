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
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/services"
	"github.com/pkg/errors"
)

// SendHTMLEmail sends validation email to new
func SendHTMLEmail(to []string, user *models.User, message string, subject string, resetPassword bool) error {
	fromEmail := configuration.GetResp().EmailNotification.Email
	from := configuration.GetResp().EmailNotification.From
	password := configuration.GetResp().EmailNotification.Password
	smtpHost := configuration.GetResp().EmailNotification.SMTPHost
	smtpPort := configuration.GetResp().EmailNotification.SMTPPort
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	var token string
	var err error
	var link string

	if !resetPassword {
		token, err = GenerateValidationToken(uint64(user.ID), configuration.GetResp().Token.ExpiryTime.ValidationToken)
		if err != nil {
			log.Println("Validation Token Generation error")
			return errors.Wrap(err, "Unable to generate validation token")
		}
		link = fmt.Sprintf("http://%s:%s/api/v1/auth/token/%s", configuration.GetResp().Server.Domain, configuration.GetResp().Server.Port, token)
	} else {
		token, err = services.CreateToken(user)
		if err != nil {
			log.Println("Reset Token Generation error")
			return errors.Wrap(err, "Unable to generate reset token")
		}
		link = fmt.Sprintf("%s%s", constants.ResetPasswordPath, token)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Unable to get working directory")
	}

	for ; strings.Split(cwd, "/")[len(strings.Split(cwd, "/"))-1] != "notifications-microservice"; cwd = filepath.Dir(cwd) {
	}

	t, err := template.ParseFiles(fmt.Sprintf("%s/shared/template/email.html", cwd))
	if err != nil {
		log.Println("Template File can't be opened")
		return errors.Wrap(err, "Unable to open template file")
	}
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	err = t.Execute(&body, struct {
		Link    string
		Message string
	}{
		Link:    link,
		Message: message,
	})
	if err != nil {
		log.Println("Unable to write to template")
		return errors.Wrap(err, "Unable to write to template")
	}

	// Authentication
	auth := smtp.PlainAuth("", fromEmail, password, addr)

	//  Sending email.
	_ = auth
	err = smtp.SendMail(addr, nil, from, to, body.Bytes())
	if err != nil {
		log.Println("Unable to send email", err.Error())
		return errors.Wrap(err, "Unable to send email")
	}
	return nil
}
