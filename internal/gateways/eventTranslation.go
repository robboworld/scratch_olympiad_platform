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

type EventTranslationGateway interface {
	CreateEventTranslation(eventTranslation models.EventTranslationCore) (newEventTranslation models.EventTranslationCore, err error)
	UpdateEventTranslation(eventTranslation models.EventTranslationCore) (updatedEventTranslation models.EventTranslationCore, err error)
	DeleteEventTranslation(id uint) (err error)
	// GetEventTranslationByEventId TODO: Когда будет несколько переводов, переделать под получения списка переводов
	GetEventTranslationByEventId(eventId uint) (eventTranslation models.EventTranslationCore, err error)
	GetEventTranslationById(id uint) (eventTranslation models.EventTranslationCore, err error)
	DoesExistEventTranslation(eventId uint) (exists bool, err error)
}

type EventTranslationGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (e EventTranslationGatewayImpl) CreateEventTranslation(
	eventTranslation models.EventTranslationCore,
) (newEventTranslation models.EventTranslationCore, err error) {
	if err = e.postgresClient.Db.Create(&eventTranslation).Clauses(clause.Returning{}).Error; err != nil {
		return models.EventTranslationCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return eventTranslation, nil
}

func (e EventTranslationGatewayImpl) UpdateEventTranslation(eventTranslation models.EventTranslationCore) (
	updatedEventTranslation models.EventTranslationCore, err error,
) {
	updateStruct := map[string]interface{}{
		"name":        eventTranslation.Name,
		"description": eventTranslation.Description,
	}

	if err = e.postgresClient.Db.Model(&eventTranslation).Clauses(clause.Returning{}).Updates(updateStruct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.EventTranslationCore{}, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return models.EventTranslationCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	if err = e.postgresClient.Db.First(&eventTranslation, eventTranslation.ID).Error; err != nil {
		return models.EventTranslationCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return eventTranslation, nil
}

func (e EventTranslationGatewayImpl) DeleteEventTranslation(id uint) (err error) {
	result := e.postgresClient.Db.Where("id = ?", id).Delete(&models.EventTranslationCore{})
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

func (e EventTranslationGatewayImpl) GetEventTranslationByEventId(eventId uint) (
	eventTranslation models.EventTranslationCore, err error,
) {
	if err = e.postgresClient.Db.Where("event_id = ?", eventId).Take(&eventTranslation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return eventTranslation, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return eventTranslation, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return eventTranslation, nil
}

func (e EventTranslationGatewayImpl) GetEventTranslationById(id uint) (
	eventTranslation models.EventTranslationCore, err error,
) {
	if err = e.postgresClient.Db.Where("id = ?", id).Take(&eventTranslation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return eventTranslation, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return eventTranslation, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return eventTranslation, nil
}

func (e EventTranslationGatewayImpl) DoesExistEventTranslation(eventId uint) (exists bool, err error) {
	if err = e.postgresClient.Db.Where("event_id = ?", eventId).
		Take(&models.EventTranslationCore{}).Error; err != nil {
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
