package integration_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"users-service/config"
	"users-service/internal/dto"
	"users-service/internal/logger"
	"users-service/internal/model"
	_ "users-service/internal/model"
	"users-service/internal/rabbitmq"
	"users-service/internal/redis"
	"users-service/internal/repository"
	"users-service/internal/router"
	"users-service/internal/validate"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	redisLib "github.com/redis/go-redis/v9"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

var app *fiber.App


type fiberRoundTripper struct {
	app *fiber.App
}

var emailJhonDoe string



func (rt fiberRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.app.Test(req)
}

var users []dto.CreateUserFromJsonFileDTO
var DB *gorm.DB

func TestMain(m *testing.M) {
	err := config.SetupEnvVar()
	if err != nil {
		log.Panic(err.Error())
	}
	if err = logger.StartLogger(); err != nil {
		log.Panic(err.Error())
	}

	if err := rabbitmq.ConnectToRabbitMQ(); err != nil {
		logger.ZapLogger.Error(err.Error())
	}
	
	validate.StartValidator()
	db, err := repository.ConnectToDatabase()
	if err != nil {
		log.Panic(err.Error())
	}
	redisClient, err := redis.ConnectToRedis()
	if err != nil {
		log.Panic(err.Error())
	}

	
	
	DB = db
	app = router.SetupApp()
	sqldb, err := db.DB()
	if err != nil {
		log.Panic(err.Error())
	}
	
	if err := cleanDatabases(db, redisClient, false); err != nil {
		log.Panic(err.Error())
	}
	code := m.Run()
	if err := cleanDatabases(db, redisClient, true); err != nil {
		log.Panic(err.Error())
	}
	sqldb.Close()
	redisClient.Close()
	
	os.Exit(code)
}

func newExpect(t *testing.T) *httpexpect.Expect {
	client := &http.Client{
		Transport: fiberRoundTripper{app: app},
	}
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost",
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
	})
}

func cleanDatabases(db *gorm.DB, redisClient *redisLib.Client, isTheEnd bool) error {
	
 	projectRoot := config.FindProjectRoot()
	envPath := filepath.Join(projectRoot, "config", "users.json")
	data, err := os.ReadFile(envPath)
	if err != nil {
		return  err
	}
	if err := json.Unmarshal(data, &users); err != nil {
		return  err
	}
	emails := funk.Map(users, func(users dto.CreateUserFromJsonFileDTO) string {
		return  users.Email
	}).([]string)
	if !isTheEnd {
		logger.ZapLogger.Info("All of users that weren't in users.json were deleted" )
		db.Where("email NOT IN ?", emails).Delete(&model.UserModel{})
	} else {
		logger.ZapLogger.Info("All of users were deleted" )
		db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.UserModel{})
	}
	redisClient.FlushDB(context.Background())
	return  nil
}

func TestMessage(t *testing.T) {
	e := newExpect(t)
	e.GET("/"). 
	Expect().
	Status(200).JSON(). 
	Object()
}


