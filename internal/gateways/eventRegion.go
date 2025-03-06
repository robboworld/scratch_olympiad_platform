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

type EventRegionGateway interface {
	SetRegionForEvent(eventId, regionId uint) error
	UnsetRegionForEvent(eventId, regionId uint) error

	DoesExistEventRegion(eventId, regionId uint) (bool, error)
}

type EventRegionGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (e EventRegionGatewayImpl) SetRegionForEvent(eventId, regionId uint) error {
	relation := models.EventRegionCore{
		EventID:  eventId,
		RegionID: regionId,
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
			Message: consts.ErrRegionAlreadyAssigned,
		}
	}
	return nil
}

func (e EventRegionGatewayImpl) UnsetRegionForEvent(eventId, regionId uint) error {
	result := e.postgresClient.Db.Where(
		"event_id = ? AND region_id = ?",
		eventId,
		regionId,
	).Delete(&models.EventRegionCore{})
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

func (e EventRegionGatewayImpl) DoesExistEventRegion(eventId, regionId uint) (bool, error) {
	if err := e.postgresClient.Db.Where("event_id = ? AND region_id = ?", eventId, regionId).
		Take(&models.EventRegionCore{}).Error; err != nil {
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
