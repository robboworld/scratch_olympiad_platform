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

type CountryGateway interface {
	GetCountryById(id uint) (country models.CountryCore, err error)
	GetAllCountries(offset, limit int) (countries []models.CountryCore, countRows uint, err error)
}

type CountryGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (c CountryGatewayImpl) GetCountryById(id uint) (country models.CountryCore, err error) {
	if err = c.postgresClient.Db.Where("id = ?", id).Take(&country).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return country, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrCountryNotFoundInDB,
			}
		}
		return country, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return country, nil
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
