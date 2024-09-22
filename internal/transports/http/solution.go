package http

import (
	"github.com/gin-gonic/gin"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type SolutionHandler struct {
	loggers logger.Loggers
}

func (h SolutionHandler) SetupSolutionRoutes(router *gin.Engine) {
	solutionGroup := router.Group("/solution")
	{
		solutionGroup.POST("/", h.UploadSolution)
		solutionGroup.GET("/download", h.DownloadSolution)
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

	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if fileHeader.Size >= consts.MaxSolutionFileSize {
		h.loggers.Err.Printf("%s", "File is too large. The maximum allowed size is 1 GB.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is too large. The maximum allowed size is 1 GB."})
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

	tempFile, err := os.CreateTemp("./internal/tmp_upload", "upload-*."+filenameParts[1])
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temporary file"})
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
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

	uploadedFilename := filepath.Base(tempFile.Name())
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"filename": uploadedFilename})
}

func (h SolutionHandler) DownloadSolution(c *gin.Context) {
	filename := c.Query("filename")
	if filename == "" {
		h.loggers.Err.Printf("%s", "filename not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename not provided"})
		return
	}
	filePath := "./internal/tmp_upload/" + filename

	role := c.Value(consts.KeyRole).(models.Role)
	accessRoles := []models.Role{models.RoleStudent, models.RoleSuperAdmin}
	if !utils.DoesHaveRole(role, accessRoles) {
		h.loggers.Err.Printf("%s", consts.ErrAccessDenied)
		c.JSON(http.StatusForbidden, gin.H{"error": consts.ErrAccessDenied})
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		h.loggers.Err.Printf("%s", "file not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
		return
	}
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Content-Type", "mp4/sb3")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Writer.Write(buf)
	c.Status(http.StatusOK)
	c.File(filename)
}
