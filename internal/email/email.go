package email

import (
	"fmt"
	"os"
	"users-service/internal/dto"
	"users-service/internal/logger"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

var (
	gomailD *gomail.Dialer
	email string
	password string
	s gomail.SendCloser
)


const (
	host = "smtp.gmail.com"
	port = 587
)

var emailQueue = make(chan dto.SendEmailDTO, 100)
func StartEmailWorker() {
	go func() {
		for e := range emailQueue {
			m := gomail.NewMessage()
			m.SetHeader("From", email)
			m.SetHeader("To", e.To)
			m.SetHeader("Subject", e.Subject)
			m.SetBody("text/plain", e.Text)

			if err := gomail.Send(s, m); err != nil {
				logger.ZapLogger.Error("error sending email", zap.Error(err))
			} else {
				logger.ZapLogger.Info("email was sent", zap.String("to", e.To))
			}
		}
	}()
}

func InitGomail() error {
	email = os.Getenv("SERVICE_EMAIL")
	password = os.Getenv("SERVICE_PASSWORD")
	if email == "" || password == "" {
		logger.ZapLogger.Error("service email or service password is empty")
		return  fmt.Errorf("service email or service password is empty")
	}
	gomailD = gomail.NewDialer(host, port, email, password)
	logger.ZapLogger.Info("gomail is ready!")
	var err error
	s, err = gomailD.Dial()
	if err != nil {
		return  err
	}
	return  nil
}

func SendEmail(input dto.SendEmailDTO) error {
	if input.Subject == "" {
		input.Subject = "Email From NewsPaper"
	}
	if input.Text == "" || input.To == "" {
		logger.ZapLogger.Error("text or to is empty")
		return  fmt.Errorf("text or to is empty")
	}
	

	

	emailQueue <- input
	logger.ZapLogger.Info("new input joined queue")

	return  nil
}