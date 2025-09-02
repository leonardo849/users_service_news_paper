package email

import (
	"fmt"
	"os"
	"users-service/internal/logger"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

var (
	gomailD *gomail.Dialer
	email string
	password string
)

const (
	host = "smtp.gmail.com"
	port = 587
)

func InitGomail() error {
	email = os.Getenv("SERVICE_EMAIL")
	password = os.Getenv("SERVICE_PASSWORD")
	if email == "" || password == "" {
		logger.ZapLogger.Error("service email or service password is empty")
		return  fmt.Errorf("service email or service password is empty")
	}
	gomailD = gomail.NewDialer(host, port, email, password)
	logger.ZapLogger.Info("gomail is ready!")
	return  nil
}

func SendEmail(to string, subject string, text string) error {
	if subject == "" {
		subject = "Email From NewsPaper"
	}
	if text == "" || to == "" {
		logger.ZapLogger.Error("text or to is empty")
		return  fmt.Errorf("text or to is empty")
	}
	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", text)

	if err := gomailD.DialAndSend(m); err != nil {
		logger.ZapLogger.Error("error in sending email", zap.Error(err))
		return  err
	}
	return  nil
}