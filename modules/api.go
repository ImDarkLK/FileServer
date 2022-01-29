package modules

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

var FileIds = make(map[string]string)

func DownloadFile(c echo.Context) error {
	link := c.Request().URL.Path

	if strings.Contains(link, "dl/name") {
		fileId := c.Param("name")
		filePath := fmt.Sprintf("downloads/%s", fileId)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return c.String(http.StatusNotFound, "File not found")
		}
		return c.Attachment(filePath, fileId)
	}

	fileId := c.Param("id")
	for id, name := range FileIds {
		if fileId == id {
			filePath := fmt.Sprintf("downloads/%s", name)
			return c.Attachment(filePath, name)
		}
	}
	return c.String(http.StatusNotFound, "File not found")
}

func HandleUpload(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	log.Println("Upload:", file.Filename)
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("downloads/" + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	downloadId := RandomString(6)
	urlPath := "http://" + c.Request().Host + "/dl"
	shortLink := fmt.Sprintf("<a href=\"/dl/id/%s\">%s</a>", downloadId, urlPath+"/id/"+downloadId)
	longLink := fmt.Sprintf("<a href=\"/dl/name/%s\">%s</a>", file.Filename, urlPath+"/name/"+file.Filename)
	FileIds[downloadId] = file.Filename

	return c.HTML(http.StatusOK, fmt.Sprintf("<h2>Your file uploaded!</h2><p>File name: %s<br>Download Links:<br>	%s (short link)<br>		%s (long link)</p>", file.Filename, shortLink, longLink))
}
