package services

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"github.com/spf13/viper"
	"github.com/thanhpk/randstr"
	"net/http"
	"time"
)

type Tokens struct {
	Access  string
	Refresh string
}

type UserClaims struct {
	jwt.StandardClaims
	Id   uint
	Role models.Role
}

type AuthService interface {
	SignUp(newUser models.UserCore) error
	SignIn(email, password string) (Tokens, error)
	Refresh(token string) (string, error)
	ConfirmActivation(link string) (Tokens, error)
	ForgotPassword(email string) error
	ResetPassword(resetLink string) error
}

type AuthServiceImpl struct {
	userGateway     gateways.UserGateway
	countryGateway  gateways.CountryGateway
	settingsGateway gateways.SettingsGateway
}

func (a AuthServiceImpl) ConfirmActivation(link string) (Tokens, error) {
	activationByLink, err := a.settingsGateway.GetActivationByLink()
	if err != nil {
		return Tokens{Access: "", Refresh: ""}, err
	}
	if !activationByLink {
		return Tokens{Access: "", Refresh: ""}, utils.ResponseError{
			Code:    http.StatusServiceUnavailable,
			Message: consts.ErrActivationLinkUnavailable,
		}
	}
	activationLinkHash := utils.GetHashString(link)
	user, err := a.userGateway.GetUserByActivationLink(activationLinkHash)
	if err != nil {
		return Tokens{Access: "", Refresh: ""}, err
	}
	if err = a.userGateway.SetIsActive(user.ID, true); err != nil {
		return Tokens{Access: "", Refresh: ""}, err
	}
	access, err := generateToken(user, viper.GetDuration("auth_access_token_ttl"), []byte(viper.GetString("auth_access_signing_key")))
	if err != nil {
		return Tokens{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	refresh, err := generateToken(user, viper.GetDuration("auth_refresh_token_ttl"), []byte(viper.GetString("auth_refresh_signing_key")))
	if err != nil {
		return Tokens{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return Tokens{Access: access, Refresh: refresh}, nil
}

func (a AuthServiceImpl) Refresh(token string) (string, error) {
	claims, err := parseToken(token, []byte(viper.GetString("auth_refresh_signing_key")))
	if err != nil {
		return "", err
	}
	user := models.UserCore{
		ID:   claims.Id,
		Role: claims.Role,
	}
	newAccessToken, err := generateToken(user, viper.GetDuration("auth_access_token_ttl"), []byte(viper.GetString("auth_access_signing_key")))
	if err != nil {
		return "", utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return newAccessToken, nil
}

func (a AuthServiceImpl) SignIn(email, password string) (Tokens, error) {
	user, err := a.userGateway.GetUserByEmail(email)
	if err != nil {
		return Tokens{}, err
	}
	if err = utils.ComparePassword(user.Password, password); err != nil {
		return Tokens{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrIncorrectPasswordOrEmail,
		}
	}
	if !user.IsActive {
		return Tokens{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrUserIsNotActive,
		}
	}
	access, err := generateToken(user, viper.GetDuration("auth_access_token_ttl"), []byte(viper.GetString("auth_access_signing_key")))
	if err != nil {
		return Tokens{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	refresh, err := generateToken(user, viper.GetDuration("auth_refresh_token_ttl"), []byte(viper.GetString("auth_refresh_signing_key")))
	if err != nil {
		return Tokens{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return Tokens{Access: access, Refresh: refresh}, nil
}

func (a AuthServiceImpl) SignUp(newUser models.UserCore) error {
	if !utils.IsValidEmail(newUser.Email) {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrIncorrectPasswordOrEmail,
		}
	}
	exist, err := a.userGateway.DoesExistEmail(0, newUser.Email)
	if err != nil {
		return err
	}
	if exist {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrEmailAlreadyInUse,
		}
	}
	if len(newUser.Password) < 8 {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrShortPassword,
		}
	}
	exist, err = a.countryGateway.DoesExistName(0, newUser.Country)
	if err != nil {
		return err
	}
	if !exist {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrCountryNotFoundInDB,
		}
	}
	activationLink := randstr.String(20)
	activationLinkHash := utils.GetHashString(activationLink)
	passwordHash := utils.HashPassword(newUser.Password)
	newUser.Password = passwordHash
	newUser.ActivationLink = activationLinkHash
	newUser, err = a.userGateway.CreateUser(newUser)
	if err != nil {
		return err
	}

	activationByLink, err := a.settingsGateway.GetActivationByLink()
	if err != nil {
		return err
	}
	var subject, body string
	if activationByLink {
		subject = "Ваша ссылка активации аккаунта"
		body = "<p>Перейдите по ссылке " + viper.GetString("activation_path") +
			activationLink + " для активации вашего аккаунта.</p>"
	} else {
		subject = "Активация аккаунта"
		body = "<p>На данный момент активация по ссылке недоступна. Ждите активации от администратора.</p>"
	}
	if err = utils.SendEmail(subject, newUser.Email, body); err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

func (a AuthServiceImpl) ForgotPassword(email string) error {
	user, err := a.userGateway.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if !user.IsActive {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrUserIsNotActive,
		}
	}

	resetPasswordLink := randstr.String(20)
	resetPasswordLinkHash := utils.GetHashString(resetPasswordLink)
	err = a.userGateway.SetPasswordResetLink(user.ID, resetPasswordLinkHash)
	if err != nil {
		return err
	}

	subject := "Ваша ссылка на сброс пароля (действует " + viper.GetString("auth_password_reset_link_at") + " минут)"
	body := "<p>Ссылка для сброса пароля: " + viper.GetString("reset_password_path") + resetPasswordLink + "</p>"
	if err = utils.SendEmail(subject, user.Email, body); err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

func (a AuthServiceImpl) ResetPassword(resetLink string) error {
	resetLinkHash := utils.GetHashString(resetLink)
	user, err := a.userGateway.GetUserByPasswordResetLink(resetLinkHash)
	if err != nil {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrPasswordResetLinkInvalid,
		}
	}
	if user.PasswordResetLinkAt.Before(time.Now()) {
		return utils.ResponseError{
			Code:    http.StatusGone,
			Message: consts.ErrPasswordResetLinkExpired,
		}
	}
	newPassword := randstr.String(8)
	newPasswordHash := utils.HashPassword(newPassword)
	err = a.userGateway.SetPassword(user.ID, newPasswordHash)
	if err != nil {
		return err
	}
	err = a.userGateway.SetPasswordResetLink(user.ID, "")
	if err != nil {
		return err
	}

	subject := "Ваш новый пароль"
	body := "<p>Новый пароль: " + newPassword + "</p>"
	if err = utils.SendEmail(subject, user.Email, body); err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

func generateToken(user models.UserCore, duration time.Duration, signingKey []byte) (token string, err error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(duration * time.Second)),
		},
		Id:   user.ID,
		Role: user.Role,
	}
	ss := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = ss.SignedString(signingKey)
	return token, err
}

func parseToken(token string, key []byte) (*UserClaims, error) {
	data, err := jwt.ParseWithClaims(token, &UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
	claims, ok := data.Claims.(*UserClaims)
	if err != nil {
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			return &UserClaims{}, utils.ResponseError{
				Code:    http.StatusUnauthorized,
				Message: consts.ErrTokenExpired,
			}
		}
		return &UserClaims{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	if !ok {
		return &UserClaims{}, utils.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: consts.ErrNotStandardToken,
		}
	}
	return claims, nil
}
