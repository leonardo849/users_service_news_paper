package integration_test

import (
	"os"
	"testing"
	"users-service/internal/dto"
	"users-service/internal/helper"
)

const userTestName string = "TestUserDev@"
var token string
var input = dto.CreateUserDTO{
	Username: userTestName,
	Email: "Batman123@gmail.com",
	Password: ";u8zG9D4b$",
	Fullname: "User Test Dev 123",
}

var idInput string 

func TestCreateUser(t *testing.T) {
	e := newExpect(t)
	res:= e.POST("/users/create"). 
	WithJSON(map[string]string{
		"username": input.Username,
		"email": input.Email,
		"password": input.Password,
		"fullname": input.Fullname,
	}).Expect(). 
	Status(201).JSON().Object()
	idInput = res.Value("id").String().NotEmpty().Raw()
	

	e.POST("/users/create"). 
	WithJSON(map[string]string{
		"username": input.Username,
		"email": input.Email,
		"password": input.Password,
		"fullname": input.Fullname,
	}).Expect(). 
	Status(409)
}

var tokenJhonDoe string

func TestLoginUser(t *testing.T) {
	
	e := newExpect(t)
	res := e.POST("/users/login"). 
	WithJSON(map[string]string{
		"email": input.Email,
		"password": input.Password,
	}). 
	Expect().
	Status(200).JSON().Object()

	token = res.Value("token").String().NotEmpty().Raw()

	e.POST("/users/login"). 
	WithJSON(map[string]string{
		"email": input.Email,
		"password": input.Password + "432423423",
	}). 
	Expect().
	Status(401)

	e.POST("/users/login"). 
	WithJSON(map[string]string{
		"email": input.Email + "4333",
		"password": input.Password + "432423423",
	}). 
	Expect().
	Status(404)


	chosenEmail := os.Getenv("EMAIL_JHONDOE")
	chosenPassword := os.Getenv("PASSWORD_JHONDOE")
	res = e.POST("/users/login"). 
	WithJSON(map[string]string{
		"email": chosenEmail,
		"password": chosenPassword,
	}). 
	Expect(). 
	Status(200).JSON().Object()

	tokenJhonDoe = res.Value("token").String().NotEmpty().Raw()
}

func TestFindOneUser(t *testing.T) {
	e := newExpect(t)
	e.GET("/users/one/" + idInput).
    WithHeader("Authorization", "Bearer " + token).
    Expect().
    Status(200).
    JSON().Object()
}

func TestUpdateOneUser(t *testing.T) {
	e := newExpect(t)
	e.PUT("/users/update/" + idInput). 
	WithJSON(map[string]string{
		"username": input.Username + "3234",
	}). 
	WithHeader("Authorization", "Bearer " + token).
	Expect().
	Status(200). 
	JSON().Object()
}


func TestUpdateOneUserRole(t *testing.T) {
	e := newExpect(t)
	e.PATCH("/users/update/role/" + idInput). 
	WithJSON(map[string]string{
		"role": helper.Developer,
	}). 
	WithHeader("Authorization", "Bearer " + tokenJhonDoe). 
	Expect(). 
	Status(200). 
	JSON().Object()
}

func TestMetrics(t *testing.T) {
	e := newExpect(t)
	e.GET("/metrics"). 
	WithHeader("Authorization", "Bearer " + tokenJhonDoe). 
	Expect().
	Status(200)
}