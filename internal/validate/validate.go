package validate

import (
	"unicode"
	"users-service/internal/helper"
	"users-service/internal/logger"

	"github.com/go-playground/validator/v10"
	"github.com/thoas/go-funk"
)

var Validate *validator.Validate


func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasMinLen  = len(password) >= 8
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func checkRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	return  funk.Contains(helper.Roles, role)
}

func StartValidator() {
	Validate = validator.New()
	Validate.RegisterValidation("strongpassword", strongPassword)
	Validate.RegisterValidation("role", checkRole)
	logger.ZapLogger.Info("Var Validate is ready!")
}