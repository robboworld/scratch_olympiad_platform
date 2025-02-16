package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
)

type RegionService interface {
	GetRegionsByCountryId(countryId uint, page, pageSize *int) (regions []models.RegionCore, countRows uint, err error)
}

type RegionServiceImpl struct {
	regionGateway gateways.RegionGateway
}

func (r RegionServiceImpl) GetRegionsByCountryId(
	countryId uint,
	page,
	pageSize *int,
) (regions []models.RegionCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return r.regionGateway.GetRegionsByCountryId(countryId, offset, limit)
}
