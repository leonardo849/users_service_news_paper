package integration_test

import (
	"testing"
	"users-service/internal/dto"
)

const userTestName string = "TestUser1@"
var token string

var input = dto.CreateUserDTO{
	Username: userTestName,
	Email: "Batman123@gmail.com",
	Password: ";u8zG9D4b$",
	Fullname: "User Test 123",
}

func TestCreateUser(t *testing.T) {
	e := newExpect(t)
	e.POST("/users/create"). 
	WithJSON(map[string]string{
		"username": input.Username,
		"email": input.Email,
		"password": input.Password,
		"fullname": input.Fullname,
	}).Expect(). 
	Status(201)

	e.POST("/users/create"). 
	WithJSON(map[string]string{
		"username": input.Username,
		"email": input.Email,
		"password": input.Password,
		"fullname": input.Fullname,
	}).Expect(). 
	Status(409)
}

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
}