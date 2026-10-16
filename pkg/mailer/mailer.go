package mailer

import (
	"fmt"
	"net"
	"net/smtp"
)

type Mailer interface {
	SendActivationCode(toEmail, code string) error
}

type SMTPMailer struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPMailer(host, port, username, password, from string) *SMTPMailer {
	return &SMTPMailer{Host: host, Port: port, Username: username, Password: password, From: from}
}

func (m *SMTPMailer) SendActivationCode(toEmail, code string) error {
	subject := "Код подтверждения регистрации"
	body := fmt.Sprintf("Ваш код подтверждения: %s\n\nКод действителен 15 минут.", code)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n%s",
		m.From, toEmail, subject, body,
	)

	var auth smtp.Auth
	if m.Username != "" {
		auth = smtp.PlainAuth("", m.Username, m.Password, m.Host)
	}

	addr := net.JoinHostPort(m.Host, m.Port)

	return smtp.SendMail(addr, auth, m.From, []string{toEmail}, []byte(msg))
}
