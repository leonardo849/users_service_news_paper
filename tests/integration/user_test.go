package integration_test

import (
	"testing"
	"users-service/internal/dto"
)

const userTestName string = "TestUser1@"

func TestCreateUser(t *testing.T) {
	input := dto.CreateUserDTO{
		Username: userTestName,
		Email: "Batman123@gmail.com",
		Password: ";u8zG9D4b$",
		Fullname: "User Test 123",
	}
	e := newExpect(t)
	e.POST("/users/create"). 
	WithJSON(map[string]string{
		"username": input.Username,
		"email": input.Email,
		"password": input.Password,
		"fullname": input.Fullname,
	}).Expect(). 
	Status(201)
}