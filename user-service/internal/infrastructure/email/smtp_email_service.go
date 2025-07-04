package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
)

type SMTPEmailService struct {
	from     string
	host     string
	port     int
	username string
	password string
}

func NewSMTPEmailService() *SMTPEmailService {
	return &SMTPEmailService{
		from:     os.Getenv("SMTP_FROM"),
		host:     os.Getenv("SMTP_HOST"),
		port:     getEnvAsInt("SMTP_PORT", 587),
		username: os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASSWORD"),
	}
}

func (s *SMTPEmailService) SendWelcomeEmail(to string) error {
	subject := "Добро пожаловать!"
	body := os.Getenv("WELCOME_EMAIL_TEMPLATE")

	if s.username == "" || s.password == "" || s.from == "" || s.host == "" {
		return fmt.Errorf("SMTP config is not fully set")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.host, s.port, s.username, s.password)
	fmt.Printf("Sending email to %s via %s:%d\n", to, s.host, s.port)
	return d.DialAndSend(m)
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	var val int
	fmt.Sscanf(valStr, "%d", &val)
	return val
}
