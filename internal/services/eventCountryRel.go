package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
)

type EventCountryService interface {
	CreateRel(eventId, countryId uint) (models.EventCountryRelCore, error)
	DeleteRel(eventId, countryId uint) error
}

type EventCountryServiceImpl struct {
	eventCountryGateway gateways.EventCountryRelGateway
}

func (e EventCountryServiceImpl) CreateRel(eventId, countryId uint) (models.EventCountryRelCore, error) {
	return e.eventCountryGateway.CreateRel(models.EventCountryRelCore{EventID: eventId, CountryID: countryId})
}

func (e EventCountryServiceImpl) DeleteRel(eventId, countryId uint) error {
	return e.eventCountryGateway.DeleteRel(models.EventCountryRelCore{EventID: eventId, CountryID: countryId})
}
