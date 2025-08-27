package rabbitmq

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"users-service/internal/logger"

	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var isRabbitMQon = false


type clientI interface {
	Publish(queue string, message string) error
	CloseRabbit()
}

type client struct {
	conn *amqp.Connection
	ch *amqp.Channel
}



var rabbitClient *client

func ConnectToRabbitMQ() (error) {
	isRabbitOn := strings.ToLower(os.Getenv("RABBIT_ON"))
	if isRabbitOn == "true" {
		isRabbitMQon = true 
	} 
	logger.ZapLogger.Info("Is rabbit going to be on?", zap.Bool("rabbit value", isRabbitMQon))
	if !isRabbitMQon {
		return  fmt.Errorf("rabbit mq is going to be off")
	}

	
	
	uriRabbit := os.Getenv("RABBIT_URI")
	if uriRabbit == "" {
		err := fmt.Errorf("rabbit_uri is empty")
		logger.ZapLogger.Error("error in get rabbit uri", zap.Error(err))
		return err
	}
	var conn *amqp.Connection
	const maxTries = 11
	secondDelay := os.Getenv("SECOND_DELAY")
	var secondInt int
	var err error
	if secondDelay == "" {
		secondInt = 1
	} else {
		secondInt, err = strconv.Atoi(secondDelay)
		if err != nil {
			secondInt = 1
		} 
	}
	for i := 0; i < maxTries; i++ {
		conn, err = amqp.Dial(uriRabbit)
		if err != nil {
			logger.ZapLogger.Error("error in Connect To rabbit mq", zap.Error(err))
			time.Sleep(time.Duration(secondInt) * time.Second)
			if i == maxTries - 1 {
				return  err
			}
		} else {
			break
		}

	}
	
	logger.ZapLogger.Info("conn is estabilished")
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		logger.ZapLogger.Info("closing conn")
	 	logger.ZapLogger.Error("error in get channel", zap.Error(err))
	 	return  err
	}
	logger.ZapLogger.Info("channel is ready")

	client := &client{
		conn: conn,
		ch: ch,
	}
	rabbitClient = client
	return nil
}

func (c *client) Publish(queue string, message string) error {
	if c.ch == nil {
		logger.ZapLogger.Error("channel doesn't exist")
		return fmt.Errorf("channel doesn't exist")
	}
	return  c.ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		},
	)
}


func (c *client) CloseRabbit() {
	c.ch.Close()
	c.conn.Close()
	logger.ZapLogger.Info("closing rabbitmq")
}



func GetRabbitMQClient() clientI {
	if isRabbitMQon && rabbitClient != nil {
		return  rabbitClient
	} else {
		return  &fakeClient{}
	}
}