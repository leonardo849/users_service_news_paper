package validate

import (
	"template-backend/internal/logger"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func StartValidator() {
	Validate = validator.New()
	logger.ZapLogger.Info("Var Validate is ready!")
}