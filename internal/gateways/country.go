package gateways

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type CountryGateway interface {
	GetAllCountries(offset, limit int) (countries []models.CountryCore, countRows uint, err error)
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
