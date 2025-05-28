package server

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/internal/services"

	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware(errLogger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(consts.AuthHeader)
		if authHeader == "" {
			c.Set(consts.KeyId, uint(0))
			c.Set(consts.KeyRole, models.RoleAnonymous)
			c.Next()
			return
		}

		if err := validateAuthHeader(authHeader); err != nil {
			errLogger.Printf("%s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
		}

		headerParts := strings.Split(authHeader, " ")
		userId, userRole, err := getUserFromAuthentication(headerParts[1])
		if err != nil {
			errLogger.Printf("%s", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}

		c.Set(consts.KeyId, userId)
		c.Set(consts.KeyRole, userRole)
		c.Next()
	}
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func getUserFromAuthentication(token string) (id uint, role models.Role, err error) {
	data, err := jwt.ParseWithClaims(token, &services.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("auth_access_signing_key")), nil
		})
	if data == nil {
		return 0, models.RoleAnonymous, errors.New(consts.ErrEmptyDataWithClaims)
	}

	claims, ok := data.Claims.(*services.UserClaims)
	if !ok {
		return 0, models.RoleAnonymous, errors.New(consts.ErrNotStandardToken)
	}
	if err != nil {
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			return 0, models.RoleAnonymous, errors.New(consts.ErrTokenExpired)
		}
		return 0, models.RoleAnonymous, err
	}

	return claims.Id, claims.Role, nil
}

func validateAuthHeader(authHeader string) error {
	// headerParts should be = ["Bearer", "<accessToken>"]
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return errors.New("invalid authorization header format")
	}
	if headerParts[0] != "Bearer" {
		return errors.New("invalid authorization header format")
	}

	return nil
}
