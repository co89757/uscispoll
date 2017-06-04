package main

import (
	"encoding/json"
	"io"
	"log"

	"gopkg.in/gomail.v2"
)

//Mailer contains info on email sender and receiver
type Mailer struct {
	SenderEmail    string   `json:"senderEmail"`
	Password       string   `json:"password"`
	SMTPServerHost string   `json:"smtpServerHost"`
	SMTPServerPort int      `json:"smtpServerPort"`
	ReceiverEmails []string `json:"receiverEmails"`
}

// NewMailer inits a new mailer based on configuration
func NewMailer(in io.Reader) (mail *Mailer, err error) {
	decoder := json.NewDecoder(in)
	mail = new(Mailer)
	err = decoder.Decode(mail)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

//SendEmail sends the mail with subject and body set
func (mailer Mailer) SendEmail(subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailer.SenderEmail)
	m.SetHeader("To", mailer.ReceiverEmails...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(mailer.SMTPServerHost, mailer.SMTPServerPort, mailer.SenderEmail, mailer.Password)
	if err := dialer.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}
	return nil
}
