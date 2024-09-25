package http

import (
	"github.com/gin-gonic/gin"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/internal/services"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type SolutionHandler struct {
	loggers         logger.Loggers
	solutionService services.SolutionService
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

	link, err := h.solutionService.CreateSolution(fileHeader)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"file_link": link})
}

func (h SolutionHandler) DownloadSolution(c *gin.Context) {
	filename := c.Query("filename")
	if filename == "" {
		h.loggers.Err.Printf("%s", "filename not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename not provided"})
		return
	}

	accessLink := c.Query("access_link")
	if accessLink == "" {
		h.loggers.Err.Printf("%s", "access link not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "access link not provided"})
		return
	}

	buf, err := h.solutionService.DownloadSolution(filename, accessLink)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "mp4/sb3")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Writer.Write(buf)
	c.Status(http.StatusOK)
}
