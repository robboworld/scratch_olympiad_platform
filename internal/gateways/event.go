package gateways

import (
	"errors"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type EventGateway interface {
	CreateEvent(event models.EventCore) (newEvent models.EventCore, err error)
	UpdateEvent(event models.EventCore) (updatedEvent models.EventCore, err error)
	GetEventById(id uint) (event models.EventCore, err error)
	GetAllEvents(offset, limit int) (events []models.EventCore, countRows uint, err error)
	GetEventsByCountryId(countryId uint, offset, limit int) (events []models.EventCore, countRows uint, err error)
	GetEventsByCountryIdAndRegionId(countryId uint, regionId uint, offset, limit int) (events []models.EventCore, countRows uint, err error)
}

type EventGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (e EventGatewayImpl) CreateEvent(event models.EventCore) (newEvent models.EventCore, err error) {
	if err = e.postgresClient.Db.Create(&event).Clauses(clause.Returning{}).Error; err != nil {
		return models.EventCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return event, nil
}

func (e EventGatewayImpl) UpdateEvent(event models.EventCore) (updatedEvent models.EventCore, err error) {
	if err := e.postgresClient.Db.Model(&event).Clauses(clause.Returning{}).
		Take(&models.EventCore{}, event.ID).
		Updates(map[string]interface{}{
			"name":        event.Name,
			"description": event.Description,
			"start_date":  event.StartDate,
			"end_date":    event.EndDate,
		}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.EventCore{}, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrEventNotFoundInDB,
			}
		}
		return models.EventCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return event, nil
}

func (e EventGatewayImpl) GetEventById(id uint) (event models.EventCore, err error) {
	if err = e.postgresClient.Db.First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.EventCore{}, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrEventNotFoundInDB,
			}
		}
		return models.EventCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return event, nil
}

func (e EventGatewayImpl) GetAllEvents(offset, limit int) (events []models.EventCore, countRows uint, err error) {
	query := e.postgresClient.Db.Model(&models.EventCore{})

	var count int64
	result := query.Count(&count)
	if result.Error != nil {
		return nil, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}

	result = query.Limit(limit).Offset(offset).Find(&events)
	if result.Error != nil {
		return []models.EventCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	return events, uint(count), result.Error
}

func (e EventGatewayImpl) GetEventsByCountryIdAndRegionId(
	countryID, regionID uint,
	offset, limit int,
) (events []models.EventCore, countRows uint, err error) {
	query := e.postgresClient.Db.Model(&models.EventCore{}).
		Joins("JOIN event_country_cores ec ON ec.event_id = event_cores.id").
		Joins("JOIN event_region_cores er ON er.event_id = event_cores.id").
		Where("ec.country_id = ? AND er.region_id = ?", countryID, regionID)

	var count int64
	result := query.Count(&count)
	if result.Error != nil {
		return []models.EventCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}

	result = query.Limit(limit).Offset(offset).Find(&events)
	if result.Error != nil {
		return nil, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	return events, uint(count), nil
}

func (e EventGatewayImpl) GetEventsByCountryId(
	countryID uint,
	offset, limit int,
) (events []models.EventCore, countRows uint, err error) {
	query := e.postgresClient.Db.Model(&models.EventCore{}).
		Joins("JOIN event_country_cores ec ON ec.event_id = event_cores.id").
		Where("ec.country_id = ?", countryID)

	var count int64
	result := query.Count(&count)
	if result.Error != nil {
		return []models.EventCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}

	result = query.Limit(limit).Offset(offset).Find(&events)
	if result.Error != nil {
		return nil, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	return events, uint(count), nil
}
