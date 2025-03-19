package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type EventRegionService interface {
	CreateRel(eventId, regionId uint) (models.EventRegionRelCore, error)
	DeleteRel(eventId, regionId uint) error
}

type EventRegionServiceImpl struct {
	regionGateway       gateways.RegionGateway
	eventGateway        gateways.EventGateway
	eventCountryGateway gateways.EventCountryRelGateway
	eventRegionGateway  gateways.EventRegionRelGateway
}

func (e EventRegionServiceImpl) CreateRel(eventId, regionId uint) (models.EventRegionRelCore, error) {
	region, err := e.regionGateway.GetRegionById(regionId)
	if err != nil {
		return models.EventRegionRelCore{}, err
	}

	// Проверка доступа мероприятия стране, которой принадлежит регион
	exist, err := e.eventCountryGateway.DoesExistRel(models.EventCountryRelCore{EventID: eventId, CountryID: region.CountryID})
	if err != nil {
		return models.EventRegionRelCore{}, err
	}
	if !exist {
		return models.EventRegionRelCore{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrRegionCountryNotLinkedToEvent,
		}
	}

	return e.eventRegionGateway.CreateRel(models.EventRegionRelCore{EventID: eventId, RegionID: regionId})
}

func (e EventRegionServiceImpl) DeleteRel(eventId, regionId uint) error {
	return e.eventRegionGateway.DeleteRel(models.EventRegionRelCore{EventID: eventId, RegionID: regionId})
}
