package gateways

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type NominationGateway interface {
	GetAllNominations(offset, limit int) (nominations []models.NominationCore, countRows uint, err error)
}

type NominationGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (n NominationGatewayImpl) GetAllNominations(offset, limit int) (nominations []models.NominationCore, countRows uint, err error) {
	var count int64
	result := n.postgresClient.Db.Limit(limit).Offset(offset).Find(&nominations)
	if result.Error != nil {
		return []models.NominationCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return nominations, uint(count), result.Error
}
