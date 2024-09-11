package utils

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
	"net/smtp"
	"time"
)

func SendEmail(subject, to, body string) (err error) {
	from := viper.GetString("mail_username")
	pass := viper.GetString("mail_password")
	e := email.NewEmail()
	e.From = "Robbo <" + from + ">"
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(body)

	auth := smtp.PlainAuth("", from, pass, viper.GetString("smtp_server_host"))
	return e.Send(viper.GetString("smtp_server_address"), auth)
}

func HashPassword(s string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(hashed)
}

func ComparePassword(hashed string, normal string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(normal))
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GetOffsetAndLimit(page, pageSize *int) (offset, limit int) {
	if page == nil || pageSize == nil {
		limit = -1
		offset = 0
	} else {
		offset = (*page - 1) * *pageSize
		limit = *pageSize
	}
	return
}

func DoesHaveRole(clientRole models.Role, roles []models.Role) bool {
	for _, role := range roles {
		if role.String() == clientRole.String() {
			return true
		}
	}
	return false
}

func GetHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func StringPointerToString(p *string) string {
	var s string
	if p != nil {
		s = *p
	}
	return s
}

func BoolPointerToBool(p *bool) bool {
	var b bool
	if p != nil {
		b = *p
	}
	return b
}

func CalculateUserAge(birthdate time.Time) int {
	today := time.Now()
	age := today.Year() - birthdate.Year()

	if today.YearDay() < birthdate.YearDay() {
		age--
	}
	return age
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		return nil, ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Gin context not found in request context",
		}
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		return nil, ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Gin context has an invalid type",
		}
	}
	return gc, nil
}
