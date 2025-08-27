package integration_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"users-service/config"
	"users-service/internal/logger"
	"users-service/internal/model"
	_ "users-service/internal/model"
	"users-service/internal/redis"
	"users-service/internal/repository"
	"users-service/internal/router"
	"users-service/internal/validate"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	redisLib "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var app *fiber.App


type fiberRoundTripper struct {
	app *fiber.App
}





func (rt fiberRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.app.Test(req)
}

var DB *gorm.DB

func TestMain(m *testing.M) {
	err := config.SetupEnvVar()
	if err != nil {
		log.Panic(err.Error())
	}
	if err = logger.StartLogger(); err != nil {
		log.Panic(err.Error())
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
	
	cleanDatabases(db, redisClient)
	code := m.Run()
	cleanDatabases(db, redisClient)
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

func cleanDatabases(db *gorm.DB, redisClient *redisLib.Client) {
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.UserModel{})
	redisClient.FlushDB(context.Background())
}

func TestMessage(t *testing.T) {
	e := newExpect(t)
	e.GET("/"). 
	Expect().
	Status(200).JSON(). 
	Object()

}