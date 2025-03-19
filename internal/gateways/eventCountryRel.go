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

type EventCountryRelGateway interface {
	CreateRel(rel models.EventCountryRelCore) (models.EventCountryRelCore, error)
	DeleteRel(rel models.EventCountryRelCore) error

	DoesExistRel(rel models.EventCountryRelCore) (bool, error)
	GetEventsByCountryId(countryId uint, offset, limit int) (events []models.EventCore, countRows uint, err error)
	GetEventsByCountryIdAndRegionId(countryID, regionID uint, offset, limit int) (events []models.EventCore, countRows uint, err error)
}

type EventCountryRelGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (e EventCountryRelGatewayImpl) CreateRel(rel models.EventCountryRelCore) (models.EventCountryRelCore, error) {
	if err := e.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		event := models.EventCore{ID: rel.EventID}
		country := models.CountryCore{ID: rel.CountryID}
		if err := tx.First(&event).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.First(&country).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.Where("event_id = ? AND country_id = ?", rel.EventID, rel.CountryID).
			First(&models.EventCountryRelCore{}).Error; err == nil {
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
		return models.EventCountryRelCore{}, err
	}
	return rel, nil
}

func (e EventCountryRelGatewayImpl) DeleteRel(rel models.EventCountryRelCore) error {
	if err := e.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where(
			"event_id = ? AND country_id = ?",
			rel.EventID, rel.CountryID,
		).Delete(&models.EventCountryRelCore{})
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

func (e EventCountryRelGatewayImpl) DoesExistRel(rel models.EventCountryRelCore) (bool, error) {
	if err := e.postgresClient.Db.Where("event_id = ? AND country_id = ?", rel.EventID, rel.CountryID).
		Take(&models.EventCountryRelCore{}).Error; err != nil {
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
