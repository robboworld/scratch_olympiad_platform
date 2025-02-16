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
	GetEventById(id uint) (event models.EventCore, err error)
	GetAllEvents(offset, limit int) (events []models.EventCore, countRows uint, err error)
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

func (e EventGatewayImpl) GetEventById(id uint) (event models.EventCore, err error) {
	if err = e.postgresClient.Db.First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.EventCore{}, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
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
	var count int64
	result := e.postgresClient.Db.Limit(limit).Offset(offset).Find(&events)
	if result.Error != nil {
		return []models.EventCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return events, uint(count), result.Error
}
