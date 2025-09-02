package integration_test

import (
	"testing"
	"users-service/internal/dto"
	"users-service/internal/helper"

	"log"
	"github.com/thoas/go-funk"
)

const userTestName string = "TestUserDev@"
var tokenDev string
var inputDev = dto.CreateUserDTO{
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
		"username": inputDev.Username,
		"email": inputDev.Email,
		"password": inputDev.Password,
		"fullname": inputDev.Fullname,
	}).Expect(). 
	Status(201).JSON().Object()
	idInput = res.Value("id").String().NotEmpty().Raw()
	

	e.POST("/users/create"). 
	WithJSON(map[string]string{
		"username": inputDev.Username,
		"email": inputDev.Email,
		"password": inputDev.Password,
		"fullname": inputDev.Fullname,
	}).Expect(). 
	Status(409)
}






func TestLoginUser(t *testing.T) {
	
	e := newExpect(t)
	res := e.POST("/users/login"). 
	WithJSON(map[string]string{
		"email": inputDev.Email,
		"password": inputDev.Password,
	}). 
	Expect().
	Status(200).JSON().Object()

	tokenDev = res.Value("token").String().NotEmpty().Raw()

	e.POST("/users/login"). 
	WithJSON(map[string]string{
		"email": inputDev.Email,
		"password": inputDev.Password + "432423423",
	}). 
	Expect().
	Status(401)

	e.POST("/users/login"). 
	WithJSON(map[string]string{
		"email": inputDev.Email + "4333",
		"password": inputDev.Password + "432423423",
	}). 
	Expect().
	Status(404)

	
}


var ceo map[string]interface{} = map[string]interface{}{
	"ceoFromJson": nil,
	"token": nil,
}

func TestLoginCeo(t *testing.T) {
	
	ceo["ceoFromJson"] = funk.Find(users, func(user dto.CreateUserFromJsonFileDTO) bool {
		return user.Role == helper.Ceo 
	}).(dto.CreateUserFromJsonFileDTO)

	var ceoFromJson dto.CreateUserFromJsonFileDTO
	var ok bool
	if ceoFromJson, ok = ceo["ceoFromJson"].(dto.CreateUserFromJsonFileDTO); !ok {
		t.Error("error in get ceo from ceo map")
	} 
	e := newExpect(t)
	resp := e.POST("/users/login"). 
	WithJSON(map[string]string{
		"email": ceoFromJson.Email,
		"password": ceoFromJson.Password,
	}). 
	Expect(). 
	Status(200).JSON().Object() 
	
	ceo["token"] = resp.Value("token").String().Raw()
}

var dev map[string]interface{} = map[string]interface{}{
	"devFromJson": nil,
	"token": nil,
}

func TestLoginDev(t *testing.T) {
	dev["devFromJson"] = funk.Find(users, func(user dto.CreateUserFromJsonFileDTO) bool {
		return user.Role == helper.Developer
	}).(dto.CreateUserFromJsonFileDTO)
	devFromJson, ok := dev["devFromJson"].(dto.CreateUserFromJsonFileDTO)
	if !ok {
		t.Error("error in dev[devFromJson] to dto.CreateUserFromJsonFileDTO")
	}
	e := newExpect(t)
	resp := e.POST("/users/login"). 
	WithJSON(map[string]string{
		"email": devFromJson.Email,
		"password": devFromJson.Password,
	}). 
	Expect(). 
	Status(200).JSON().Object()

	dev["token"] = resp.Value("token").String().Raw()
}

var journalist map[string]interface{} = map[string]interface{}{
	"journalistFromJson": nil,
	"token": nil,
}

func TestLoginJournalist(t *testing.T) {
	journalist["journalistFromJson"] = funk.Find(users, func(user dto.CreateUserFromJsonFileDTO) bool {
		return user.Role == helper.Journalist
	}).(dto.CreateUserFromJsonFileDTO)
	journalistFromJson, ok := journalist["journalistFromJson"].(dto.CreateUserFromJsonFileDTO)
	if !ok {
		t.Error("error in journalist[journalistFromJson] to dto.CreateUserFromJsonFileDTO")
	}
	log.Print(journalistFromJson.Email)
	e := newExpect(t)
	resp := e.POST("/users/login"). 
	WithJSON(map[string]string{
		"email": journalistFromJson.Email,
		"password": journalistFromJson.Password,
	}). 
	Expect(). 
	Status(200).JSON().Object()

	dev["token"] = resp.Value("token").String().Raw()
}

func TestFindOneUser(t *testing.T) {
	e := newExpect(t)
	e.GET("/users/one/" + idInput).
    WithHeader("Authorization", "Bearer " + tokenDev).
    Expect().
    Status(200).
    JSON().Object()
}

func TestUpdateOneUser(t *testing.T) {
	e := newExpect(t)
	e.PUT("/users/update/" + idInput). 
	WithJSON(map[string]string{
		"username": inputDev.Username + "3234",
	}). 
	WithHeader("Authorization", "Bearer " + tokenDev).
	Expect().
	Status(200). 
	JSON().Object()
}


func TestUpdateOneUserRole(t *testing.T) {
	tokenCeo, ok := ceo["token"].(string)
	if !ok {
		t.Error("error in ceo token to string")
	}
	e := newExpect(t)
	e.PATCH("/users/update/role/" + idInput). 
	WithJSON(map[string]string{
		"role": helper.Developer,
	}).
	WithHeader("Authorization", "Bearer " + tokenCeo).
	Expect().
	Status(200). 
	JSON().Object()

	
}

func TestMetrics(t *testing.T) {
	tokenCeo, ok := ceo["token"].(string)
	if !ok {
		t.Error("error in ceo token to string")
	}
	e := newExpect(t)
	e.GET("/metrics"). 
	WithHeader("Authorization", "Bearer " + tokenCeo). 
	Expect(). 
	Status(200)
}