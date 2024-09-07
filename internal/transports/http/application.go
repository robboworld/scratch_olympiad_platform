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

type ApplicationHandler struct {
	loggers            logger.Loggers
	applicationService services.ApplicationService
}

func (h ApplicationHandler) SetupApplicationRoutes(router *gin.Engine) {
	applicationGroup := router.Group("/application")
	{
		applicationGroup.POST("/", h.CreateApplication)
	}
}

func (h ApplicationHandler) CreateApplication(c *gin.Context) {
	var input models.NewApplication
	if err := c.ShouldBindJSON(&input); err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.Value(consts.KeyId).(uint)
	role := c.Value(consts.KeyRole).(models.Role)
	roleStudent := models.RoleStudent
	accessRoles := []*models.Role{&roleStudent}
	if !utils.DoesHaveRole(role, accessRoles) {
		h.loggers.Err.Printf("%s", consts.ErrAccessDenied)
		c.JSON(http.StatusForbidden, gin.H{"error": consts.ErrAccessDenied})
		return
	}

	application := models.ApplicationCore{
		AuthorID:   userID,
		Nomination: input.Nomination,
		Link:       utils.StringPointerToString(input.Link),
		Note:       utils.StringPointerToString(input.Note),
	}
	newApplication, err := h.applicationService.CreateApplication(application)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	applicationHttp := models.ApplicationHTTP{}
	applicationHttp.FromCore(newApplication)
	c.JSON(http.StatusOK, gin.H{"application": applicationHttp})
}
