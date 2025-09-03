package helper

import (
	"fmt"
	"time"
	"users-service/config"
	"users-service/internal/logger"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)



func GenerateJWT(id string, updatedAt time.Time,  role string) (string, error){
	if config.Key == "" {
		logger.ZapLogger.Error("secret word jwt is empty", zap.String("function", "generate jwt"))
		return  "", fmt.Errorf("secret word jwt is empty")
	}
	secretKeyBytes := []byte(config.Key)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id": id,
		"updatedAt": updatedAt,
		"role": role,
		"exp": time.Now().Add((time.Hour * 24) * 30).Unix(),
    })

	tokenString, err := token.SignedString(secretKeyBytes)
	if err != nil {
		logger.ZapLogger.Error("error in token signed string", zap.Error(err), zap.String("function", "generate jwt"))
		return  "", err
	}

	return  tokenString, nil
}