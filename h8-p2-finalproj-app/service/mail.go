package service

import (
	"errors"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendMail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "mrdrummerman123@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	host := os.Getenv("SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	if err != nil {
		return errors.New("invalid SMTP port")
	}
	d := gomail.NewDialer(host, port, user, pass)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
