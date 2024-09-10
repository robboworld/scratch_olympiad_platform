package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
)

type CountryService interface {
	GetAllCountries(page, pageSize *int) (countries []models.CountryCore, countRows uint, err error)
}

type CountryServiceImpl struct {
	countryGateway gateways.CountryGateway
}

func (c CountryServiceImpl) GetAllCountries(page, pageSize *int) (countries []models.CountryCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return c.countryGateway.GetAllCountries(offset, limit)
}
