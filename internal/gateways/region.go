package gateways

import (
	"errors"
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"gorm.io/gorm"
	"net/http"
)

type RegionGateway interface {
	GetAllRegions(offset, limit int) (regions []models.RegionCore, countRows uint, err error)
	GetRegionsByCountryId(countryId uint, offset, limit int) (regions []models.RegionCore, countRows uint, err error)
	DoesExistRegion(countryId, id uint, name string) (bool, error)
}

type RegionGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (r RegionGatewayImpl) GetAllRegions(offset, limit int) (regions []models.RegionCore, countRows uint, err error) {
	var count int64
	result := r.postgresClient.Db.Limit(limit).Offset(offset).Find(&regions)
	if result.Error != nil {
		return []models.RegionCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return regions, uint(count), result.Error
}

func (r RegionGatewayImpl) GetRegionsByCountryId(countryId uint, offset, limit int) (regions []models.RegionCore, countRows uint, err error) {
	var count int64
	result := r.postgresClient.Db.Limit(limit).Offset(offset).Where("country_id = ?", countryId).
		Find(&regions)
	if result.Error != nil {
		return []models.RegionCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return regions, uint(count), result.Error
}

func (r RegionGatewayImpl) DoesExistRegion(countryId, id uint, name string) (bool, error) {
	if err := r.postgresClient.Db.Where("id != ? AND name = ? AND country_id = ?", id, name, countryId).
		Take(&models.RegionCore{}).Error; err != nil {
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
