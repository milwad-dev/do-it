package services

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

func SendEmail(data struct {
	Subject string
	Name    string
	Body    string
}, templatePath, toDest string) error {
	// Email auth info
	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASSWORD")
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")

	// Parse HTML template
	tmpl, err := template.ParseFiles("internal/templates/" + templatePath)
	if err != nil {
		// TODO: Add log
		return err
	}

	var body bytes.Buffer

	// Get app data
	appName := os.Getenv("APP_NAME")
	appMail := os.Getenv("APP_MAIL")

	// Write email headers
	body.WriteString("MIME-Version: 1.0\r\n")
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	body.WriteString(fmt.Sprintf("From: %s\r\n", appMail))
	body.WriteString(fmt.Sprintf("To: %s\r\n", toDest))
	body.WriteString(fmt.Sprintf("Subject: %s - %s\r\n", appName, data.Subject))
	body.WriteString("\r\n")

	// Execute the template with data
	err = tmpl.Execute(&body, data)
	if err != nil {
		// TODO: Add log
		return err
	}

	// Set up auth
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email
	err = smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{toDest},
		body.Bytes(),
	)

	if err != nil {
		// TODO: Add log
		return err
	}

	return nil
}
