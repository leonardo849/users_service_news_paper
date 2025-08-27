package rabbitmq

import "users-service/internal/logger"

type fakeClient struct {
}

func (c *fakeClient) Publish(queue string, message string) error {
	log := "[fake] publish: queue " + queue + "message " + message
	logger.ZapLogger.Info(log)
	return nil
}

func (c *fakeClient) CloseRabbit() {
	logger.ZapLogger.Info("[fake] closing rabbit")
}

