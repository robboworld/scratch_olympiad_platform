package http

import (
	"github.com/gin-gonic/gin"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/internal/services"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
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

	newUser := models.UserCore{
		Email:          input.Email,
		Password:       input.Password,
		Firstname:      input.Firstname,
		Lastname:       input.Lastname,
		Middlename:     utils.StringPointerToString(input.Middlename),
		Nickname:       input.Nickname,
		Role:           models.RoleStudent,
		IsActive:       false,
		ActivationLink: utils.GetHashString(time.Now().String()),
	}

	err := h.authService.SignUp(newUser)
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
