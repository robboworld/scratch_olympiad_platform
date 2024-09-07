package http

import (
	"github.com/gin-gonic/gin"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/internal/services"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
	"io"
	"net/http"
	"strconv"
)

type ProjectHandler struct {
	loggers        logger.Loggers
	projectService services.ProjectService
}

func (h ProjectHandler) SetupProjectRoutes(router *gin.Engine) {
	projectGroup := router.Group("/project")
	{
		projectGroup.GET("/", h.GetProjectById)
		projectGroup.POST("/", h.UpdateProject)
	}
}

func (h ProjectHandler) GetProjectById(c *gin.Context) {
	projectId := c.Query("id")
	atoi, err := strconv.Atoi(projectId)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrAtoi})
		return
	}

	userID := c.Value(consts.KeyId).(uint)
	role := c.Value(consts.KeyRole).(models.Role)
	project, err := h.projectService.GetProjectById(
		uint(atoi),
		userID,
		role,
	)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project.Json})
}

func (h ProjectHandler) UpdateProject(c *gin.Context) {
	projectId := c.Query("id")
	atoi, err := strconv.Atoi(projectId)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrAtoi})
		return
	}

	dataBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect json body"})
		return
	}

	project := models.ProjectCore{
		ID:   uint(atoi),
		Json: string(dataBytes),
	}

	_, err = h.projectService.UpdateProject(project)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": project.ID})
}
