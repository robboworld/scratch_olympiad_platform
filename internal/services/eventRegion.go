package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type EventRegionService interface {
	SetRegionForEvent(eventId, regionId uint, clientRole models.Role) error
	UnsetRegionForEvent(eventId, regionId uint, clientRole models.Role) error
}

type EventRegionServiceImpl struct {
	regionGateway       gateways.RegionGateway
	eventGateway        gateways.EventGateway
	eventCountryGateway gateways.EventCountryGateway
	eventRegionGateway  gateways.EventRegionGateway
}

func (e EventRegionServiceImpl) SetRegionForEvent(eventId, regionId uint, clientRole models.Role) error {
	// Нет доступа для User
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	region, err := e.regionGateway.GetRegionById(regionId)
	if err != nil {
		return err
	}

	_, err = e.eventGateway.GetEventById(eventId)
	if err != nil {
		return err
	}

	// Проверка доступа мероприятия стране, которой принадлежит регион
	exist, err := e.eventCountryGateway.DoesExistEventCountry(eventId, region.CountryID)
	if err != nil {
		return err
	}
	if !exist {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrRegionCountryNotLinkedToEvent,
		}
	}

	return e.eventRegionGateway.SetRegionForEvent(eventId, regionId)
}

func (e EventRegionServiceImpl) UnsetRegionForEvent(eventId, regionId uint, clientRole models.Role) error {
	// Нет доступа для User
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}
	return e.eventRegionGateway.UnsetRegionForEvent(eventId, regionId)
}
