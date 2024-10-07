package http

import (
	"github.com/gin-gonic/gin"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type AvatarHandler struct {
	loggers logger.Loggers
}

func (h AvatarHandler) SetupAvatarRoutes(router *gin.Engine) {
	avatarGroup := router.Group("/avatar")
	{
		avatarGroup.POST("/", h.UploadAvatar)
	}
}

func (h AvatarHandler) UploadAvatar(c *gin.Context) {
	// TODO rm to avatar service
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	filename := strings.Split(fileHeader.Filename, ".")
	if len(filename) != 2 {
		h.loggers.Err.Printf("%s", "incorrect filename")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect filename"})
		return
	}

	if filename[1] != "png" && filename[1] != "jpg" {
		h.loggers.Err.Printf("%s", "incorrect file format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect file format"})
		return
	}

	tempFile, err := os.CreateTemp("./internal/tmp_upload", "upload-*."+filename[1])
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temporary file"})
		return
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	tempFile.Write(fileBytes)

	err = os.Chmod(tempFile.Name(), 0644)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to chmod file"})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"filename": filepath.Base(tempFile.Name())})
}
