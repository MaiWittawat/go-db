package storer

import (
	"fmt"
	"go-rebuild/internal/handler"
	"go-rebuild/internal/storer"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type storerHandler struct {
	service storer.Storer
}

func NewStorerHandler(service storer.Storer) handler.StorerHandler {
	return &storerHandler{service: service}
}

func (h *storerHandler) GetFileUrl(c *gin.Context) {
	fileName := c.Query("object")
	url, err := h.service.GetFileUrl(c.Request.Context(), fileName, "GET")
	if err != nil {
		log.Printf("Error dowload file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to dowload file from storer: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "get url success", "url": url})
}

func (h *storerHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to get file from form: %v", err)})
		return
	}

	fileReader, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to open uploaded file: %v", err)})
		return
	}
	defer fileReader.Close()

	if err := h.service.Upload(
		c.Request.Context(),
		fileReader,
		fileHeader.Size,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
	); err != nil {
		log.Printf("Error during file upload: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file to storer: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "upload success"})
}
