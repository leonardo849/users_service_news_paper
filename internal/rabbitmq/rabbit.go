package rabbitmq

import (
	"fmt"
	"os"
	"users-service/internal/logger"

	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var isRabbitMQon = false


type client struct {
	conn *amqp.Connection
	ch *amqp.Channel
}

func (c *client) CloseRabbit() {
	c.ch.Close()
	c.conn.Close()
	logger.ZapLogger.Info("closing rabbitmq")
}

var RabbitClient *client

func ConnectToRabbitMQ() (*client, error) {
	isRabbitOn := strings.ToLower(os.Getenv("RABBIT_ON"))
	if isRabbitOn == "true" {
		isRabbitMQon = true 
	} 
	logger.ZapLogger.Info("Is rabbit going to be on?", zap.Bool("rabbit value", isRabbitMQon))
	if !isRabbitMQon {
		return  nil, fmt.Errorf("rabbit mq is going to be off")
	}

	
	
	uriRabbit := os.Getenv("RABBIT_URI")
	if uriRabbit == "" {
		err := fmt.Errorf("rabbit_uri is empty")
		logger.ZapLogger.Error("error in get rabbit uri", zap.Error(err))
		return nil, err
	}
	
	conn, err := amqp.Dial(uriRabbit)
	if err != nil {
		logger.ZapLogger.Error("error in Connect To rabbit mq", zap.Error(err))
		return nil, err
	}
	logger.ZapLogger.Info("conn is estabilished")
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		logger.ZapLogger.Info("closing conn")
	 	logger.ZapLogger.Error("error in get channel", zap.Error(err))
	 	return  nil, err
	}
	logger.ZapLogger.Info("channel is ready")

	client := &client{
		conn: conn,
		ch: ch,
	}
	RabbitClient = client
	return client, nil
}

