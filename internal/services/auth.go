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
	ConfirmActivation(token string) (Tokens, error)
	ForgotPassword(email string) error
	ResetPassword(token string) error
}

type AuthServiceImpl struct {
	userGateway     gateways.UserGateway
	authDataGateway gateways.AuthDataGateway
	countryGateway  gateways.CountryGateway
	regionGateway   gateways.RegionGateway
	settingsGateway gateways.SettingsGateway
}

func (a AuthServiceImpl) ConfirmActivation(token string) (Tokens, error) {
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
	activationTokenHash := utils.GetHashString(token)
	userId, err := a.authDataGateway.GetUserIdByActivationToken(activationTokenHash)
	if err != nil {
		return Tokens{Access: "", Refresh: ""}, err
	}
	user, err := a.userGateway.GetUserById(userId)
	if err != nil {
		return Tokens{Access: "", Refresh: ""}, err
	}
	if err = a.userGateway.SetIsActive(user.ID, true); err != nil {
		return Tokens{Access: "", Refresh: ""}, err
	}
	if err = a.authDataGateway.SetActivationToken(user.ID, ""); err != nil {
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
	country, err := a.countryGateway.GetCountryById(newUser.CountryID)
	if err != nil {
		return err
	}
	newUser.Country = country
	if newUser.RegionID == nil && country.HasRegions {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrRegionRequired,
		}
	}
	if newUser.RegionID != nil {
		if !country.HasRegions {
			return utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrCountryHasNoRegions,
			}
		}
		region, err := a.regionGateway.GetRegionById(*newUser.RegionID)
		if err != nil {
			return err
		}
		if newUser.CountryID != region.CountryID {
			return utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrRegionNotInCountry,
			}
		}
		newUser.Region = &region
	}
	activationToken := randstr.String(20)
	activationByLink, err := a.settingsGateway.GetActivationByLink()
	if err != nil {
		return err
	}
	var subject, body string
	if activationByLink {
		subject = "Scratch Olympiad account activation"
		body = "<p>Please follow this link to activate your Scratch Olympiad account:</p>" +
			"<p><a href='" + viper.GetString("activation_link") + activationToken + "'>" +
			viper.GetString("activation_link") + activationToken + "</a></p><br>" +
			"<p>Organizing committee of the International Scratch Creative Programming Olympiad</p>" +
			"<p><a href='mailto:scratch@creativeprogramming.org'>scratch@creativeprogramming.org</a></p>" +
			"<p><a href='https://creativeprogramming.org'>creativeprogramming.org</a></p>"
	} else {
		subject = "Scratch Olympiad account activation"
		body = "<p>Activation via the link is not available at the moment. Wait for activation from the administrator</p>" + "<br>" +
			"<p>Organizing committee of the International Scratch Creative Programming Olympiad</p>" +
			"<p>scratch@creativeprogramming.org</p>" +
			"<p>creativeprogramming.org</p>"
	}
	if err = utils.SendEmail(subject, newUser.Email, body); err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	passwordHash := utils.HashPassword(newUser.Password)
	newUser.Password = passwordHash
	user, err := a.userGateway.CreateUser(newUser)
	if err != nil {
		return err
	}

	if err = a.authDataGateway.CreateAuthData(user.ID); err != nil {
		return err
	}
	activationTokenHash := utils.GetHashString(activationToken)
	if err = a.authDataGateway.SetActivationToken(user.ID, activationTokenHash); err != nil {
		return err
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

	passwordResetToken := randstr.String(20)
	subject := "Request to reset your Scratch Olympiad account password"
	body := "<p>We have received a request to reset your account password.</p>" +
		"<p>If you did it, please follow this link (the link is active for " +
		viper.GetString("auth_password_reset_token_at") + " minutes):</p>" +
		"<p><a href='" + viper.GetString("password_reset_link") + passwordResetToken + "'>" +
		viper.GetString("password_reset_link") + passwordResetToken + "</a></p><br>" +
		"<p>If you did not do this, please just ignore this email.</p><br>" +
		"<p>Organizing committee of the International Scratch Creative Programming Olympiad</p>" +
		"<p><a href='mailto:scratch@creativeprogramming.org'>scratch@creativeprogramming.org</a></p>" +
		"<p><a href='https://creativeprogramming.org'>creativeprogramming.org</a></p>"

	if err = utils.SendEmail(subject, user.Email, body); err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	passwordResetTokenHash := utils.GetHashString(passwordResetToken)
	if err = a.authDataGateway.SetPasswordResetToken(user.ID, passwordResetTokenHash); err != nil {
		return err
	}
	return nil
}

func (a AuthServiceImpl) ResetPassword(token string) error {
	passwordResetTokenHash := utils.GetHashString(token)
	userId, err := a.authDataGateway.GetUserIdByPasswordResetToken(passwordResetTokenHash)
	if err != nil {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrPasswordResetTokenInvalid,
		}
	}

	user, err := a.userGateway.GetUserById(userId)
	if err != nil {
		return err
	}

	authData, err := a.authDataGateway.GetAuthDataByUserId(userId)
	if err != nil {
		return err
	}

	if authData.PasswordResetTokenAt.Before(time.Now()) {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrPasswordResetTokenExpired,
		}
	}
	newPassword := randstr.String(8)
	subject := "Your new Scratch Olympiad account password"
	body := "<p>Your new password:</p>" +
		"<p>" + newPassword + "</p><br>" +
		"<p>Organizing committee of the International Scratch Creative Programming Olympiad</p>" +
		"<p><a href='mailto:scratch@creativeprogramming.org'>scratch@creativeprogramming.org</a></p>" +
		"<p><a href='https://creativeprogramming.org'>creativeprogramming.org</a></p>"
	if err = utils.SendEmail(subject, user.Email, body); err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	newPasswordHash := utils.HashPassword(newPassword)
	if err = a.userGateway.SetPassword(user.ID, newPasswordHash); err != nil {
		return err
	}
	if err = a.authDataGateway.SetPasswordResetToken(user.ID, ""); err != nil {
		return err
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
