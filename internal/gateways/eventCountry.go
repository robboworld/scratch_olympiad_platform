package gateways

import (
	"errors"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"gorm.io/gorm"
	"net/http"
)

type EventCountryGateway interface {
	SetCountryForEvent(eventId, countryId uint) error
	UnsetCountryForEvent(eventId, countryId uint) error

	DoesExistEventCountry(eventId, countryId uint) (bool, error)
}

type EventCountryGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (e EventCountryGatewayImpl) SetCountryForEvent(eventId, countryId uint) error {
	relation := models.EventCountryCore{
		EventID:   eventId,
		CountryID: countryId,
	}
	result := e.postgresClient.Db.Where(relation).FirstOrCreate(&relation)
	if result.Error != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	if result.RowsAffected == 0 {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrCountryAlreadyAssigned,
		}
	}
	return nil
}

func (e EventCountryGatewayImpl) UnsetCountryForEvent(eventId, countryId uint) error {
	result := e.postgresClient.Db.Where(
		"event_id = ? AND country_id = ?",
		eventId,
		countryId,
	).Delete(&models.EventCountryCore{})
	if result.Error != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	if result.RowsAffected == 0 {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrNotFoundInDB,
		}
	}
	return nil
}

func (e EventCountryGatewayImpl) DoesExistEventCountry(eventId, countryId uint) (bool, error) {
	if err := e.postgresClient.Db.Where("event_id = ? AND country_id = ?", eventId, countryId).
		Take(&models.EventCountryCore{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return true, nil
}
