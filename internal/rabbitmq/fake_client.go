package rabbitmq

import (
	"context"
	"strings"
	"users-service/internal/logger"
	"users-service/pkg/email_dto"
	dtoSl "github.com/leonardo849/shared_library_news_paper/pkg/dto" 
)

type fakeClient struct {
}

func (c *fakeClient) PublishEmail(input email_dto.SendEmailDTO, ctx context.Context) error {
	log := "[fake] publish email. To " + strings.Join(input.To, "") + " subject "  + input.Subject + " text " + input.Text
	logger.ZapLogger.Info(log)
	return nil
}

func (c *fakeClient) CloseRabbit() {
	logger.ZapLogger.Info("[fake] closing rabbit")
}


func (c *fakeClient) createExchanges() {
	logger.ZapLogger.Info("[fake] creating exchanges")
}

func (c *fakeClient) PublishUserVerified(input dtoSl.AuthPublishUserCreated, ctx context.Context) error {
	logger.ZapLogger.Info("[fake] publishing userverified")
	return  nil
}
