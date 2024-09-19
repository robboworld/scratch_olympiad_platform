package http

import (
	"github.com/gin-gonic/gin"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"io"
	"net/http"
	"os"
	"strings"
)

const MaxFileSize = 1 * 1024 * 1024 * 1024

type SolutionHandler struct {
	loggers logger.Loggers
}

func (h SolutionHandler) SetupSolutionRoutes(router *gin.Engine) {
	solutionGroup := router.Group("/solution")
	{
		solutionGroup.POST("/", h.UploadSolution)
	}
}

func (h SolutionHandler) UploadSolution(c *gin.Context) {
	role := c.Value(consts.KeyRole).(models.Role)
	accessRoles := []models.Role{models.RoleStudent, models.RoleSuperAdmin}
	if !utils.DoesHaveRole(role, accessRoles) {
		h.loggers.Err.Printf("%s", consts.ErrAccessDenied)
		c.JSON(http.StatusForbidden, gin.H{"error": consts.ErrAccessDenied})
		return
	}

	if err := c.Request.ParseMultipartForm(1 << 30); err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	filenameParts := strings.Split(fileHeader.Filename, ".")
	if len(filenameParts) != 2 {
		h.loggers.Err.Printf("%s", "incorrect filename")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect filename"})
		return
	}

	allowedFormats := map[string]bool{"mp4": true, "sb3": true}
	if !allowedFormats[filenameParts[1]] {
		h.loggers.Err.Printf("%s", "incorrect file format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect file format"})
		return
	}

	buffer := io.LimitReader(file, MaxFileSize)
	var fileBytes []byte
	fileBytes, err = io.ReadAll(buffer)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	if len(fileBytes) >= MaxFileSize {
		h.loggers.Err.Printf("%s", "File is too large. The maximum allowed size is 1 GB.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is too large. The maximum allowed size is 1 GB."})
		return
	}

	tempFile, err := os.CreateTemp("./internal/tmp_upload", "upload-*."+filenameParts[1])
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temporary file"})
		return
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, io.LimitReader(file, MaxFileSize)); err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write file"})
		return
	}

	err = os.Chmod(tempFile.Name(), 0644)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to chmod file"})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"filename": strings.Replace(tempFile.Name(), "\\", "/", -1)})
}
