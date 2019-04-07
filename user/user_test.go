package main

import (
	"testing"

	"../user/userEntitie"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/httptest"
)

// $ go test -v

func TestNewApp(t *testing.T) {
	app := NewApp("../config.json")
	e := httptest.New(t, app)

	//GET ALL TEST
	fmt.Printf("GET ALL TEST\n\n")
	fmt.Printf("Body Response: ")
	fmt.Printf(e.GET("/user").Expect().
		Status(httptest.StatusOK).Body().Raw())
	fmt.Printf("\n\n")

	//POST TEST
	fmt.Printf("POST TEST\n\n")
	user := userEntitie.User{
		Username: "TEST_USERNAME",
		Password: "TEST_PASSWORD"}
	userToJson, _ := json.Marshal(user)

	fmt.Printf("Body: ")
	fmt.Printf(string(userToJson))

	id := e.POST("/user").WithJSON(user).Expect().
		Status(httptest.StatusOK).Body().Raw()

	fmt.Printf("\nBody Response: ")
	fmt.Printf(string(id))

	fmt.Printf("\n\n")

	//FINDBYID TEST
	fmt.Printf("FINDBYID TEST\n\n")
	var buffer bytes.Buffer
	buffer.WriteString("/user/")
	buffer.WriteString(id)

	fmt.Printf("Params: %v", id)

	fmt.Printf("\nBody Response: ")
	fmt.Printf(e.GET(buffer.String()).Expect().
		Status(httptest.StatusOK).Body().
		Contains(id).
		Contains(user.Username).
		Contains(user.Password).Raw())

	fmt.Printf("\n\n")

	//FINDBY USERNAME AND PASSWORD TEST
	//fmt.Printf("FINDBY USERNAME AND PASSWORD TEST\n\n")
	//buffer.WriteString("/user/TEST_USERNAME/TEST_PASSWORD")
	//
	//fmt.Printf("Params: TEST_USERNAME/TEST_PASSWORD")
	//
	//fmt.Printf("\nBody Response: ")
	//fmt.Printf(e.GET(buffer.String()).Expect().
	//	Status(httptest.StatusOK).Body().
	//	Contains(id).
	//	Contains(user.Username).
	//	Contains(user.Password).Raw())
	//
	//fmt.Printf("\n\n")
	//
	////UPDATE TEST
	//fmt.Printf("UPDATE TEST\n\n")
	//user.Id, _ = strconv.Atoi(id)
	//user.Username = "TEST_USERNAME_EDITED"
	//user.Password = "TEST_PASSWORD_EDITED"
	//
	//userToJson, _ = json.Marshal(user)
	//
	//fmt.Printf("Body: ")
	//fmt.Printf(string(userToJson))
	//
	//fmt.Printf("\nBody Response: ")
	//fmt.Printf(e.PUT("/user").WithJSON(user).Expect().
	//	Status(httptest.StatusOK).Body().Raw())
	//
	//fmt.Printf("\n\n")

	//FINDBYID TEST
	fmt.Printf("FINDBYID TEST\n\n")
	buffer.Reset()
	buffer.WriteString("/user/")
	buffer.WriteString(id)

	fmt.Printf("Params: %v", id)

	fmt.Printf("\nBody Response: ")
	fmt.Printf(e.GET(buffer.String()).Expect().
		Status(httptest.StatusOK).Body().
		Contains(id).
		Contains(user.Username).
		Contains(user.Password).Raw())

	fmt.Printf("\n\n")

	//DELETE TEST
	fmt.Printf("DELETE TEST\n\n")
	buffer.Reset()
	buffer.WriteString("/user/")
	buffer.WriteString(id)

	fmt.Printf("Params: %v", id)

	fmt.Printf("\nBody Response: ")
	fmt.Printf(e.DELETE(buffer.String()).Expect().
		Status(httptest.StatusOK).Body().Raw())

	fmt.Printf("\n\n")

}
