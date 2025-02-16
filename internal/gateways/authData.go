package gateways

import (
	"errors"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"time"
)

type AuthDataGateway interface {
	CreateAuthData(userId uint) (err error)
	SetPasswordResetToken(userId uint, passwordResetToken string) error
	SetActivationToken(userId uint, activationToken string) error

	GetUserIdByPasswordResetToken(passwordResetToken string) (userId uint, err error)
	GetUserIdByActivationToken(activationToken string) (userId uint, err error)
	GetAuthDataByUserId(userId uint) (authData models.AuthDataCore, err error)
}

type AuthDataGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (a AuthDataGatewayImpl) CreateAuthData(userId uint) (err error) {
	authData := models.AuthDataCore{
		ActivationToken:      "",
		PasswordResetToken:   "",
		PasswordResetTokenAt: time.Time{},
		UserID:               userId,
	}
	if err = a.postgresClient.Db.Create(&authData).Clauses(clause.Returning{}).Error; err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

func (a AuthDataGatewayImpl) SetPasswordResetToken(userId uint, passwordResetToken string) (err error) {
	var updateStruct map[string]interface{}
	if passwordResetToken == "" {
		updateStruct = map[string]interface{}{
			"password_reset_token":    "",
			"password_reset_token_at": time.Time{},
		}
	} else {
		updateStruct = map[string]interface{}{
			"password_reset_token":    passwordResetToken,
			"password_reset_token_at": time.Now().Add(time.Minute * viper.GetDuration("auth_password_reset_token_at")),
		}
	}
	if err := a.postgresClient.Db.First(&models.AuthDataCore{UserID: userId}).Updates(updateStruct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

func (a AuthDataGatewayImpl) SetActivationToken(userId uint, activationToken string) (err error) {
	var updateStruct map[string]interface{}
	if activationToken == "" {
		updateStruct = map[string]interface{}{
			"activation_token": "",
		}
	} else {
		updateStruct = map[string]interface{}{
			"activation_token": activationToken,
		}
	}
	if err := a.postgresClient.Db.First(&models.AuthDataCore{UserID: userId}).Updates(updateStruct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

func (a AuthDataGatewayImpl) GetUserIdByPasswordResetToken(passwordResetToken string) (userId uint, err error) {
	var authData models.AuthDataCore
	if err = a.postgresClient.Db.Where("password_reset_token = ?", passwordResetToken).Take(&authData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return authData.UserID, nil
}

func (a AuthDataGatewayImpl) GetUserIdByActivationToken(activationToken string) (userId uint, err error) {
	var authData models.AuthDataCore
	if err = a.postgresClient.Db.Where("activation_token = ?", activationToken).Take(&authData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return authData.UserID, nil
}

func (a AuthDataGatewayImpl) GetAuthDataByUserId(userId uint) (authData models.AuthDataCore, err error) {
	if err = a.postgresClient.Db.Where("user_id = ?", userId).Take(&authData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return authData, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return authData, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return authData, nil
}
