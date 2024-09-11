package http

import (
	"github.com/gin-gonic/gin"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/internal/services"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
	"time"
)

type AuthHandler struct {
	loggers     logger.Loggers
	authService services.AuthService
}

func (h AuthHandler) SetupAuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/sign-up", h.SignUp)
		authGroup.POST("/sign-in", h.SignIn)
		authGroup.POST("/confirm", h.Confirm)
	}
}

func (h AuthHandler) SignUp(c *gin.Context) {
	var input models.SignUp
	if err := c.ShouldBindJSON(&input); err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	birthdate, err := time.Parse(time.DateOnly, input.Birthdate)
	if err != nil {
		h.loggers.Err.Printf("%s", consts.ErrTimeParse)
		c.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrTimeParse})
		return
	}
	newUser := models.UserCore{
		Email:          input.Email,
		Password:       input.Password,
		FullName:       input.FullName,
		FullNameNative: input.FullNameNative,
		Country:        input.Country,
		City:           input.City,
		Birthdate:      birthdate,
		Role:           models.RoleStudent,
		IsActive:       false,
		ActivationLink: utils.GetHashString(time.Now().String()),
	}

	err = h.authService.SignUp(newUser)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h AuthHandler) SignIn(c *gin.Context) {
	var input models.SignIn
	if err := c.ShouldBindJSON(&input); err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authService.SignIn(input.Email, input.Password)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.Access,
		"refresh_token": tokens.Refresh,
	})
}

func (h AuthHandler) Confirm(c *gin.Context) {
	activationLink := c.Query("activationLink")
	tokens, err := h.authService.ConfirmActivation(activationLink)
	if err != nil {
		h.loggers.Err.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.Access,
		"refresh_token": tokens.Refresh,
	})
}
