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

type EventUserRelGateway interface {
	CreateRel(rel models.EventUserRelCore) (models.EventUserRelCore, error)
	DeleteRel(rel models.EventUserRelCore) error

	GetEventRoles(eventId, userId uint) ([]models.EventRole, error)
}

type EventUserRelGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (e EventUserRelGatewayImpl) CreateRel(rel models.EventUserRelCore) (models.EventUserRelCore, error) {
	if err := e.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		event := models.EventCore{ID: rel.EventID}
		user := models.UserCore{ID: rel.UserID}
		if err := tx.First(&event).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.First(&user).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.Where(
			"event_id = ? AND user_id = ? AND event_role = ?",
			rel.EventID, rel.UserID, rel.EventRole).
			First(&models.EventUserRelCore{}).Error; err == nil {
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
		if err := tx.Create(&rel).Clauses(clause.Returning{}).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		return nil
	}); err != nil {
		return models.EventUserRelCore{}, err
	}
	return rel, nil
}

func (e EventUserRelGatewayImpl) DeleteRel(rel models.EventUserRelCore) error {
	if err := e.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where(
			"event_id = ? AND user_id = ? AND event_role = ?",
			rel.EventID, rel.UserID, rel.EventRole).Delete(&models.EventUserRelCore{})
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

func (e EventUserRelGatewayImpl) GetEventRoles(eventId, userId uint) ([]models.EventRole, error) {
	var roles []models.EventRole
	if err := e.postgresClient.Db.Model(&models.EventUserRelCore{}).
		Where("event_id = ? AND user_id = ?", eventId, userId).
		Pluck("event_role", &roles).Error; err != nil {

		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return roles, nil
}
