package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"strings"
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

func validateMailConfig(mailer Mailer) error {
	if mailer.SMTPServerHost == "" || mailer.SMTPServerPort == 0 {
		return fmt.Errorf("Error: no SMTP server or port configured")
	}
	if len(mailer.ReceiverEmails) < 1 {
		return fmt.Errorf("Error: No mail receiver is specified")
	}
	if !strings.Contains(mailer.SenderEmail, "@") {
		return fmt.Errorf("Error: invalid email sender: %s", mailer.SenderEmail)
	}
	return nil

}

func composeBody(subject, body, name, from string) []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("Subject: " + subject + "\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("From: %s <%s>\r\n", name, from))
	buf.WriteString(fmt.Sprintf("Content-Type: text/html; charset=utf-8\r\n\r\n"))
	buf.WriteString(body)
	return buf.Bytes()
}

//SendEmail sends the mail with subject and body set
func (mailer Mailer) SendEmail(subject, senderName, body string) (err error) {
	err = validateMailConfig(mailer)
	if err != nil {
		return
	}
	auth := smtp.PlainAuth(
		"",
		mailer.SenderEmail,
		mailer.Password,
		mailer.SMTPServerHost)
	conn, err := smtp.Dial(fmt.Sprintf("%s:%d", mailer.SMTPServerHost, mailer.SMTPServerPort))
	if err != nil {
		return
	}
	err = conn.StartTLS(&tls.Config{ServerName: mailer.SMTPServerHost})
	if err != nil {
		return
	}
	err = conn.Auth(auth)
	if err != nil {
		return
	}
	err = conn.Mail(mailer.SenderEmail)
	if err != nil {
		if strings.Contains(err.Error(), "530 5.5.1") {
			err = fmt.Errorf("Error: Authentication failure. ")

		}
		return
	}
	for _, recv := range mailer.ReceiverEmails {
		err = conn.Rcpt(recv)
		if err != nil {
			return
		}
	}

	wc, err := conn.Data()
	if err != nil {
		return
	}
	defer wc.Close()
	bodybytes := composeBody(subject, body, senderName, mailer.SenderEmail)
	_, err = wc.Write(bodybytes)
	if err != nil {
		return
	}
	return nil
}
