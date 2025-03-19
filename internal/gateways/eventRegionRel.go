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

type EventRegionRelGateway interface {
	CreateRel(rel models.EventRegionRelCore) (models.EventRegionRelCore, error)
	DeleteRel(rel models.EventRegionRelCore) error

	DoesExistRel(rel models.EventRegionRelCore) (bool, error)
}

type EventRegionRelGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (e EventRegionRelGatewayImpl) CreateRel(rel models.EventRegionRelCore) (models.EventRegionRelCore, error) {
	if err := e.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		event := models.EventCore{ID: rel.EventID}
		region := models.RegionCore{ID: rel.RegionID}
		if err := tx.First(&event).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.First(&region).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.Where("event_id = ? AND region_id = ?", rel.EventID, rel.RegionID).
			First(&models.EventRegionRelCore{}).Error; err == nil {
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
		return models.EventRegionRelCore{}, err
	}
	return rel, nil
}

func (e EventRegionRelGatewayImpl) DeleteRel(rel models.EventRegionRelCore) error {
	if err := e.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where(
			"event_id = ? AND region_id = ?",
			rel.EventID, rel.RegionID,
		).Delete(&models.EventRegionRelCore{})
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

func (e EventRegionRelGatewayImpl) DoesExistRel(rel models.EventRegionRelCore) (bool, error) {
	if err := e.postgresClient.Db.Where("event_id = ? AND region_id = ?", rel.EventID, rel.RegionID).
		Take(&models.EventRegionRelCore{}).Error; err != nil {
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
