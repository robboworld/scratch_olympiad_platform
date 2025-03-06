package gateways

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type EventUserGateway interface {
	SetUserForEvent(eventId, userId uint, eventRole models.EventRole) error
	UnsetUserForEvent(eventId, userId uint, eventRole models.EventRole) error

	GetUserRolesForEvent(eventId, userId uint) ([]models.EventRole, error)
}

type EventUserGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (e EventUserGatewayImpl) GetUserRolesForEvent(eventId, userId uint) ([]models.EventRole, error) {
	var roles []models.EventRole
	if err := e.postgresClient.Db.Model(&models.EventUserCore{}).Select("DISTINCT(event_role)").
		Where("event_id = ? AND user_id = ?", eventId, userId).
		Pluck("event_role", &roles).Error; err != nil {

		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return roles, nil
}

func (e EventUserGatewayImpl) SetUserForEvent(eventId, userId uint, eventRole models.EventRole) error {
	relation := models.EventUserCore{
		EventID:   eventId,
		UserID:    userId,
		EventRole: eventRole,
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
			Message: consts.ErrEventRoleAlreadyAssigned,
		}
	}
	return nil
}

func (e EventUserGatewayImpl) UnsetUserForEvent(eventId, userId uint, eventRole models.EventRole) error {
	result := e.postgresClient.Db.Where(
		"event_id = ? AND user_id = ? AND event_role = ?",
		eventId,
		userId,
		eventRole,
	).Delete(&models.EventUserCore{})
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
