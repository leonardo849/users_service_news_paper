package rabbitmq

import (
	"context"
	"encoding/json"
	"users-service/internal/logger"
	"users-service/pkg/email_dto"
	dtoSl "github.com/leonardo849/shared_library_news_paper/pkg/dto" 
	constsSl "github.com/leonardo849/shared_library_news_paper/pkg/consts"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const keyEmail = constsSl.KeyEmail
const exchangeNameEmail = constsSl.ExchangeNameEmail
const exchangeNameAuthEvents = constsSl.ExchangeNameAuthEvents
const keyUserVerified = constsSl.KeyUserAuthVerified


func (c *client) createExchanges() {
	
	
	if err := c.ch.ExchangeDeclare(
		exchangeNameEmail,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logger.ZapLogger.Fatal(err.Error(), zap.String("function", "client.createexchanges"))
	}
	logger.ZapLogger.Info("exchange email_direct was declared")
	if  err := c.ch.ExchangeDeclare(
                exchangeNameAuthEvents, 
                "topic",      
                true,         
                false,        
                false,       
                false,        
                nil,          
        ); err != nil {
			logger.ZapLogger.Fatal(err.Error(), zap.String("function", "client.createexchanges"))
	}
	logger.ZapLogger.Info("exchange auth_events topic was declared")
	
}

func (c *client) PublishUserVerified(input dtoSl.AuthPublishUserCreated, ctx context.Context)  error {
	body, err := json.Marshal(input)
	if err != nil {
		logger.ZapLogger.Error("error in json marshal", zap.String("function", "client.PublishEmail"), zap.Error(err))
		return  err
	}
	err = c.ch.PublishWithContext(ctx, 
		exchangeNameAuthEvents,
		keyUserVerified,
		false, 
		false, 
		amqp091.Publishing{
			ContentType: "application/json",
			Body: body,
		},
	)
	if err != nil {
		return err
	}
	return  nil
}

func (c *client) PublishEmail(input email_dto.SendEmailDTO, ctx context.Context) error {

	body, err := json.Marshal(input)
	if err != nil {
		logger.ZapLogger.Error("error in json marshal", zap.String("function", "client.PublishEmail"), zap.Error(err))
		return  err
	}


	err = c.ch.PublishWithContext(ctx, exchangeNameEmail, keyEmail, false, false, amqp091.Publishing{
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

// func (c *client) PublishUserVerified(input )

func (c *client) CloseRabbit() {
	c.ch.Close()
	c.conn.Close()
	logger.ZapLogger.Info("rabbit was closed")
}