package server

import (
	"context"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/internal/services"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
	"time"
)

func Auth(next http.Handler, errLogger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = r.WithContext(context.WithValue(r.Context(), consts.KeyId, "0"))
			r = r.WithContext(context.WithValue(r.Context(), consts.KeyRole, models.RoleAnonymous))
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			errLogger.Printf("%s", "invalid authorization header format")
			http.Error(w, "invalid authorization header format", http.StatusBadRequest)
			return
		}
		data, err := jwt.ParseWithClaims(headerParts[1], &services.UserClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString("auth_access_signing_key")), nil
			})
		if data == nil {
			errLogger.Printf("%s", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims, ok := data.Claims.(*services.UserClaims)
		if err != nil {
			if claims.ExpiresAt.Unix() < time.Now().Unix() {
				errLogger.Printf("%s", err.Error())
				http.Error(w, consts.ErrTokenExpired, http.StatusUnauthorized)
				return
			}
			errLogger.Printf("%s", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !ok {
			errLogger.Printf("%s", consts.ErrNotStandardToken)
			http.Error(w, consts.ErrNotStandardToken, http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), consts.KeyId, claims.Id))
		r = r.WithContext(context.WithValue(r.Context(), consts.KeyRole, claims.Role))
		next.ServeHTTP(w, r)
	})
}
