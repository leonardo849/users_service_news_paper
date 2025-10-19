package integration_test

import (
	"log"
	"os"
	"testing"
	"users-service/internal/dto"
	"users-service/internal/helper"

	"github.com/thoas/go-funk"
)



const userTestName string = "TestUserDev@"
var input = dto.CreateUserDTO{
	Username: userTestName,
	Email: "",
	Password: ";u8zG9D4b$",
	Fullname: "User Test Dev 123",
}

var ceo = map[string]interface{}{
	"ceoFromJson": nil,
	"token": nil,
	"id": nil,
}

var customer = map[string]interface{}{
	"customerFromJson": nil,
	"token": nil,
	"id": nil,
}

var journalist  = map[string]interface{}{
	"journalistFromJson": nil,
	"token": nil,
	"id": nil,
}

var dev = map[string]interface{}{
	"devFromJson": nil,
	"token": nil,
	"id": nil,
}

// func findAllUsers() ([]model.UserModel, error){
// 	var users []model.UserModel
// 	result := DB.Find(&users)
// 	if result.Error != nil {
// 		return  nil, result.Error
// 	} 
// 	return  users, result.Error
// }


// func GetIds() (users []model.UserModel, err error) {
// 	users, err = findAllUsers()
// 	if err != nil {
// 		return nil, err
// 	}

// }


func TestCreateUser(t *testing.T) {
	input.Email = os.Getenv("EMAIL_INPUT")
	log.Print(input.Email)
	e := newExpect(t)
	e.POST("/users/create"). 
	WithJSON(map[string]string{
		"username": input.Username,
		"email": input.Email,
		"password": input.Password,
		"fullname": input.Fullname,
	}).Expect(). 
	Status(201).JSON().Object()


	e.POST("/users/create"). 
	WithJSON(map[string]string{
		"username": input.Username,
		"email": input.Email,
		"password": input.Password,
		"fullname": input.Fullname,
	}).Expect(). 
	Status(409)
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

func TestFindAll(t *testing.T) {
	tokenCeo := ceo["token"].(string)
	e := newExpect(t)
	resp := e.GET("/users/all").
	WithHeader("Authorization", "Bearer " + tokenCeo).
	Expect().
	Status(200).JSON().Array()

	usersI := resp.Raw()
	for _, u := range usersI {
		user := u.(map[string]interface{})
		if user["role"] == helper.Ceo && ceo["id"] == nil {
			ceo["id"] = user["id"]
		} else if user["role"] == helper.Customer && customer["id"] == nil {
			customer["id"] = user["id"]
		} else if user["role"] == helper.Developer && dev["id"] == nil {
			dev["id"] = user["id"]
		} else if journalist["id"] == nil && user["role"] == helper.Journalist{
			journalist["id"] = user["id"]
		}
	}
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



func TestLoginJournalist(t *testing.T) {
	journalist["journalistFromJson"] = funk.Find(users, func(user dto.CreateUserFromJsonFileDTO) bool {
		return user.Role == helper.Journalist
	}).(dto.CreateUserFromJsonFileDTO)
	journalistFromJson, ok := journalist["journalistFromJson"].(dto.CreateUserFromJsonFileDTO)
	if !ok {
		t.Error("error in journalist[journalistFromJson] to dto.CreateUserFromJsonFileDTO")
	}
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




func TestLoginCustomer(t *testing.T) {
	customer["customerFromJson"] = funk.Find(users, func(user dto.CreateUserFromJsonFileDTO) bool {
		return user.Role == helper.Customer
	}).(dto.CreateUserFromJsonFileDTO)
	customerFromJson, ok := customer["customerFromJson"].(dto.CreateUserFromJsonFileDTO)
	if !ok {
		t.Error("error in customer[customerFromJson] to dto.CreateUserFromJsonFileDTO")
	}
	e := newExpect(t)
	resp := e.POST("/users/login"). 
	WithJSON(map[string]string{
		"email": customerFromJson.Email,
		"password": customerFromJson.Password,
	}). 
	Expect(). 
	Status(200).JSON().Object()

	customer["token"] = resp.Value("token").String().Raw()
}




func TestFindOneUser(t *testing.T) {
	log.Print(customer["id"])
	searchedId := customer["id"].(string)
	token := ceo["token"].(string)
	e := newExpect(t)
	e.GET("/users/one/" + searchedId).
    WithHeader("Authorization", "Bearer " + token).
    Expect().
    Status(200).
    JSON().Object()
}



func TestUpdateOneUser(t *testing.T) {
	searchedId := customer["id"].(string)
	token := ceo["token"].(string)
	
	e := newExpect(t)
	e.PUT("/users/update/" + searchedId). 
	WithJSON(map[string]string{
		"username": input.Username + "3234",
	}). 
	WithHeader("Authorization", "Bearer " + token).
	Expect().
	Status(200). 
	JSON().Object()
}


func TestUpdateOneUserRole(t *testing.T) {
	searchedId := customer["id"].(string)
	tokenCeo, ok := ceo["token"].(string)
	if !ok {
		t.Error("error in ceo token to string")
	}
	e := newExpect(t)
	e.PATCH("/users/update/role/" + searchedId). 
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