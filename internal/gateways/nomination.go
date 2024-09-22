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

type NominationGateway interface {
	GetAllNominations(offset, limit int) (nominations []models.NominationCore, countRows uint, err error)
	DoesExistName(id uint, name string) (bool, error)
	GetNominationByName(name string) (nomination models.NominationCore, err error)
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

func (n NominationGatewayImpl) DoesExistName(id uint, name string) (bool, error) {
	if err := n.postgresClient.Db.Where("id != ? AND name = ?", id, name).
		Take(&models.NominationCore{}).Error; err != nil {
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

func (n NominationGatewayImpl) GetNominationByName(name string) (nomination models.NominationCore, err error) {
	if err = n.postgresClient.Db.Where("name = ?", name).Take(&nomination).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nomination, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNominationNotFoundInDB,
			}
		}
		return nomination, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nomination, nil
}
