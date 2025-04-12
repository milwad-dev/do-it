package services

import (
	"bytes"
	"fmt"
	"net/smtp"
)

func Send(data struct{ Subject string }, template string) {
	// Email auth info
	from := "b71c6c2d21b3b6"
	password := "258c9932614d3a" // TODO: Read from env
	to := "recipient@example.com"
	smtpHost := "sandbox.smtp.mailtrap.io"
	smtpPort := "587"

	// Parse HTML template
	tmpl, err := template.ParseFiles(template)
	if err != nil {
		panic(err)
	}

	var body bytes.Buffer

	// Write email headers
	body.WriteString("MIME-Version: 1.0\r\n")
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	body.WriteString(fmt.Sprintf("Subject: %s\r\n", data.Subject))
	body.WriteString("\r\n")

	// Execute the template with data
	err = tmpl.Execute(&body, data)
	if err != nil {
		panic(err)
	}

	// Set up auth
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email
	err = smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		body.Bytes(),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent successfully!")
}
