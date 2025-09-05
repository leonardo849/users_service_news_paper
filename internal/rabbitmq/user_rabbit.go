package rabbitmq

import (
	"context"
	"encoding/json"
	"users-service/internal/logger"
	"users-service/pkg/email_dto"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const keyEmail = "email"
const exchangeName = "email_direct"




func (c *client) createExchanges() {
	
	
	if err := c.ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logger.ZapLogger.Fatal(err.Error(), zap.String("function", "client.ConsumerEmail"))
	}
	logger.ZapLogger.Info("exchange email_direct was declared")

	
}



func (c *client) PublishEmail(input email_dto.SendEmailDTO, ctx context.Context) error {

	body, err := json.Marshal(input)
	if err != nil {
		logger.ZapLogger.Error("error in json marshal", zap.String("function", "client.PublishEmail"), zap.Error(err))
		return  err
	}


	err = c.ch.PublishWithContext(ctx, exchangeName, keyEmail, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body: body,
	})
	if err != nil {
		logger.ZapLogger.Error("error in publishing email", zap.String("function", "client.PublishEmail"), zap.Error(err))
		return  err
	}
	logger.ZapLogger.Info("email was sent")
	return  nil

}

func (c *client) CloseRabbit() {
	c.ch.Close()
	c.conn.Close()
	logger.ZapLogger.Info("rabbit was closed")
}