package main

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/RehanAfridikkk/word-count-Echo-API-fileupload/cmd"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Message struct {
	Upload   multipart.File `form:"upload"`
	Routines int            `form:"routines"`
}

func main() {
	e := echo.New()

	e.Use(middleware.BodyLimit("1000M"))

	e.POST("/post", postData)

	e.Logger.Fatal(e.Start(":8080"))
}

func postData(c echo.Context) error {

	var message Message
	if err := c.Request().ParseMultipartForm(5000); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to parse form"})
	}

	upload, _, err := c.Request().FormFile("upload")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing file upload"})
	}
	message.Upload = upload

	routines, err := strconv.Atoi(c.FormValue("routines"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid routines value"})
	}
	message.Routines = routines

	totalCounts, routines, timetaken, err := cmd.ProcessFile(message.Upload, message.Routines)
	timeTakenString := timetaken.String()

	return c.JSON(http.StatusOK, map[string]interface{}{"counts": totalCounts, "routines": routines, "timetaken": timeTakenString})
}
