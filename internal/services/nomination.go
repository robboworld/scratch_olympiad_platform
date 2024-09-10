package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
)

type NominationService interface {
	GetAllNominations(page, pageSize *int) (nominations []models.NominationCore, countRows uint, err error)
}

type NominationServiceImpl struct {
	nominationGateway gateways.NominationGateway
}

func (p NominationServiceImpl) GetAllNominations(page, pageSize *int) (nominations []models.NominationCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return p.nominationGateway.GetAllNominations(offset, limit)
}
