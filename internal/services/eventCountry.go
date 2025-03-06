package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type EventCountryService interface {
	SetCountryForEvent(eventId, countryId uint, clientRole models.Role) error
	UnsetCountryForEvent(eventId, countryId uint, clientRole models.Role) error
}

type EventCountryServiceImpl struct {
	countryGateway      gateways.CountryGateway
	eventGateway        gateways.EventGateway
	eventCountryGateway gateways.EventCountryGateway
}

func (e EventCountryServiceImpl) SetCountryForEvent(eventId, countryId uint, clientRole models.Role) error {
	// Нет доступа для User
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	_, err := e.countryGateway.GetCountryById(countryId)
	if err != nil {
		return err
	}

	_, err = e.eventGateway.GetEventById(eventId)
	if err != nil {
		return err
	}

	return e.eventCountryGateway.SetCountryForEvent(eventId, countryId)
}

func (e EventCountryServiceImpl) UnsetCountryForEvent(eventId, countryId uint, clientRole models.Role) error {
	// Нет доступа для User
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	return e.eventCountryGateway.UnsetCountryForEvent(eventId, countryId)
}
