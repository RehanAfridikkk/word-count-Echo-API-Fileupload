package main

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/RehanAfridikkk/word-count-Echo-API/cmd"
)

type Message struct {
	Upload   *multipart.FileHeader `form:"upload"`
	Routines int                   `form:"routines"`
}

func main() {
	e := echo.New()

	// Middleware for handling file uploads
	e.Use(middleware.BodyLimit("100M")) // Set an appropriate limit for file size

	e.POST("/post", postData)

	e.Logger.Fatal(e.Start(":8080"))
}

// Existing imports...

func postData(c echo.Context) error {
	var message Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Check if the "upload" field is present in the form data
	if message.Upload == nil {
		fmt.Println(message.Upload)
		fmt.Println(message.Routines)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing file upload"})
	}

	// Process the file content using your existing logic
	totalCounts, routines, timetaken:= cmd.ProcessFile(message.Upload.Filename, message.Routines)
	
	timeTakenString := timetaken.String()

	return c.JSON(http.StatusOK, map[string]interface{}{"counts": totalCounts, "routines": routines, "timetaken": timeTakenString})
}

