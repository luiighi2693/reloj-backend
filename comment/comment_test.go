package main

import (
	"testing"

	"../comment/commentEntitie"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/httptest"
	"strconv"
)

// $ go test -v

func TestNewApp(t *testing.T) {
	app := NewApp("../config.json")
	e := httptest.New(t, app)

	//GET ALL TEST
	fmt.Printf("GET ALL TEST\n\n")
	fmt.Printf("Body Response: ")
	fmt.Printf(e.GET("/comment").Expect().
		Status(httptest.StatusOK).Body().Raw())
	fmt.Printf("\n\n")

	//POST TEST
	fmt.Printf("POST TEST\n\n")
	comment := commentEntitie.Comment{
		Content: "TEST_CONTENT",
		Author:  "TEST_AUTHOR"}
	commentToJson, _ := json.Marshal(comment)

	fmt.Printf("Body: ")
	fmt.Printf(string(commentToJson))

	id := e.POST("/comment").WithJSON(comment).Expect().
		Status(httptest.StatusOK).Body().Raw()

	fmt.Printf("\nBody Response: ")
	fmt.Printf(string(id))

	fmt.Printf("\n\n")

	//FINDBYID TEST
	fmt.Printf("FINDBYID TEST\n\n")
	var buffer bytes.Buffer
	buffer.WriteString("/comment/")
	buffer.WriteString(id)

	fmt.Printf("Params: %v", id)

	fmt.Printf("\nBody Response: ")
	fmt.Printf(e.GET(buffer.String()).Expect().
		Status(httptest.StatusOK).Body().
		Contains(id).
		Contains(comment.Content).
		Contains(comment.Author).Raw())

	fmt.Printf("\n\n")

	//UPDATE TEST
	fmt.Printf("UPDATE TEST\n\n")
	comment.Id, _ = strconv.Atoi(id)
	comment.Content = "TEST_CONTENT_EDITED"
	comment.Author = "TEST_AUTHOR_EDITED"

	commentToJson, _ = json.Marshal(comment)

	fmt.Printf("Body: ")
	fmt.Printf(string(commentToJson))

	fmt.Printf("\nBody Response: ")
	fmt.Printf(e.PUT("/comment").WithJSON(comment).Expect().
		Status(httptest.StatusOK).Body().Raw())

	fmt.Printf("\n\n")

	//FINDBYID TEST
	fmt.Printf("FINDBYID TEST\n\n")
	buffer.Reset()
	buffer.WriteString("/comment/")
	buffer.WriteString(id)

	fmt.Printf("Params: %v", id)

	fmt.Printf("\nBody Response: ")
	fmt.Printf(e.GET(buffer.String()).Expect().
		Status(httptest.StatusOK).Body().
		Contains(id).
		Contains(comment.Content).
		Contains(comment.Author).Raw())

	fmt.Printf("\n\n")

	//DELETE TEST
	fmt.Printf("DELETE TEST\n\n")
	buffer.Reset()
	buffer.WriteString("/comment/")
	buffer.WriteString(id)

	fmt.Printf("Params: %v", id)

	fmt.Printf("\nBody Response: ")
	fmt.Printf(e.DELETE(buffer.String()).Expect().
		Status(httptest.StatusOK).Body().Raw())

	fmt.Printf("\n\n")

}
