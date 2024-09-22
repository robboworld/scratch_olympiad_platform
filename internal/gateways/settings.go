package gateways

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type SettingsGateway interface {
	SetActivationByLink(activationByCode bool) error
	GetActivationByLink() (activationByCode bool, err error)
}

type SettingsGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (s SettingsGatewayImpl) GetActivationByLink() (activationByCode bool, err error) {
	if err = s.postgresClient.Db.Model(&models.SettingsCore{}).Select("activation_by_link").Where("id = ? ", 1).
		First(&activationByCode).Error; err != nil {
		return false, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return activationByCode, nil
}

func (s SettingsGatewayImpl) SetActivationByLink(activationByCode bool) error {
	if err := s.postgresClient.Db.Model(&models.SettingsCore{ID: 1}).Updates(map[string]interface{}{
		"activation_by_link": activationByCode,
	}).Error; err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}
