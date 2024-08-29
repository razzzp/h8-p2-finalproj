package service

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail() error {
	m := gomail.NewMessage()
	m.SetHeader("From", "testmail@mail.com")
	m.SetHeader("To", "m.razif.pramuda@gmail.com")
	m.SetHeader("Subject", "Welcome to car rental app!")
	m.SetBody("text/html", "<h1>Hello, thank you for registering with us!</h1>")

	d := gomail.NewDialer("smtp.freesmtpservers.com", 25, "testmail@mail.com", os.Getenv("SMTP_PASS"))

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
