package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
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
	if err := s.postgresClient.Db.Model(&models.SettingsCore{}).Select("activation_by_link").Where("id = ? ", 1).
		First(&activationByCode).Error; err != nil {
		return false, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return activationByCode, nil
}

func (s SettingsGatewayImpl) SetActivationByLink(activationByCode bool) error {
	return s.postgresClient.Db.Model(&models.SettingsCore{ID: 1}).Updates(map[string]interface{}{
		"activation_by_link": activationByCode,
	}).Error
}
