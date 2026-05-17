package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailService struct {
	Host string
	Port string
	From string
}

func NewEmailService(host, port, from string) *EmailService {
	return &EmailService{
		Host: host,
		Port: port,
		From: from,
	}
}

func (s *EmailService) SendEmail(to string) error {
	subject := "Welcome!"
	body := "This is a test message using MailHog."

	header := make(map[string]string)
	header["From"] = s.From
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""

	var message strings.Builder
	for k, v := range header {
		fmt.Fprintf(&message, "%s: %s\r\n", k, v)
	}
	message.WriteString("\r\n")
	message.WriteString(body)

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	err := smtp.SendMail(addr, nil, s.From, []string{to}, []byte(message.String()))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
