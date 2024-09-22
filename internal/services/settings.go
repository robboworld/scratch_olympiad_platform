package services

import "github.com/robboworld/scratch_olympiad_platform/internal/gateways"

type SettingsService interface {
	SetActivationByLink(activationByCode bool) error
	GetActivationByLink() (activationByCode bool, err error)
}

type SettingsServiceImpl struct {
	settingsGateway gateways.SettingsGateway
}

func (s SettingsServiceImpl) SetActivationByLink(activationByCode bool) error {
	return s.settingsGateway.SetActivationByLink(activationByCode)
}

func (s SettingsServiceImpl) GetActivationByLink() (activationByCode bool, err error) {
	return s.settingsGateway.GetActivationByLink()
}
