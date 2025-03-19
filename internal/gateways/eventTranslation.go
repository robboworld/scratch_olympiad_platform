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
	CreateEventTranslation(eventTranslation models.EventTranslationCore) (models.EventTranslationCore, error)
	UpdateEventTranslation(eventTranslation models.EventTranslationCore) (models.EventTranslationCore, error)
	DeleteEventTranslation(id uint) error
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
) (models.EventTranslationCore, error) {
	if err := e.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		event := models.EventCore{ID: eventTranslation.EventID}
		if err := tx.First(&event).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.Where("event_id = ?", eventTranslation.EventID).
			First(&models.EventTranslationCore{}).Error; err == nil {
			return utils.ResponseError{
				Code:    http.StatusConflict,
				Message: consts.ErrRelationAlreadyExists,
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.Create(&eventTranslation).Clauses(clause.Returning{}).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		return nil
	}); err != nil {
		return models.EventTranslationCore{}, err
	}
	return eventTranslation, nil
}

func (e EventTranslationGatewayImpl) UpdateEventTranslation(eventTranslation models.EventTranslationCore) (
	models.EventTranslationCore, error,
) {
	if err := e.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		updateStruct := map[string]interface{}{
			"name":        eventTranslation.Name,
			"description": eventTranslation.Description,
		}

		if err := e.postgresClient.Db.Model(&eventTranslation).Clauses(clause.Returning{}).Updates(updateStruct).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.ResponseError{
					Code:    http.StatusBadRequest,
					Message: consts.ErrNotFoundInDB,
				}
			}
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := e.postgresClient.Db.First(&eventTranslation, eventTranslation.ID).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		return nil
	}); err != nil {
		return models.EventTranslationCore{}, err
	}
	return eventTranslation, nil
}

func (e EventTranslationGatewayImpl) DeleteEventTranslation(id uint) error {
	if err := e.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ?", id).Delete(&models.EventTranslationCore{})
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
	}); err != nil {
		return err
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
