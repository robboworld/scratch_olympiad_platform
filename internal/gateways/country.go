package gateways

import (
	"errors"
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"gorm.io/gorm"
	"net/http"
)

type CountryGateway interface {
	GetAllCountries(offset, limit int) (countries []models.CountryCore, countRows uint, err error)
	DoesExistCountry(id uint, name string) (bool, error)
}

type CountryGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (c CountryGatewayImpl) GetAllCountries(offset, limit int) (countries []models.CountryCore, countRows uint, err error) {
	var count int64
	result := c.postgresClient.Db.Limit(limit).Offset(offset).Find(&countries)
	if result.Error != nil {
		return []models.CountryCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return countries, uint(count), result.Error
}

func (c CountryGatewayImpl) DoesExistCountry(id uint, name string) (bool, error) {
	if err := c.postgresClient.Db.Where("id != ? AND name = ?", id, name).
		Take(&models.CountryCore{}).Error; err != nil {
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
