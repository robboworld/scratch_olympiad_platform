package gateways

import (
	"errors"
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type ApplicationGateway interface {
	CreateApplication(application models.ApplicationCore) (newApplication models.ApplicationCore, err error)
	DoesExistApplication(userId uint, nomination string) (bool, error)
	GetApplicationsByAuthorId(id uint, offset, limit int) (applications []models.ApplicationCore, countRows uint, err error)
	GetAllApplications(offset, limit int) (applications []models.ApplicationCore, countRows uint, err error)
}

type ApplicationGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (a ApplicationGatewayImpl) CreateApplication(application models.ApplicationCore) (newApplication models.ApplicationCore, err error) {
	if err = a.postgresClient.Db.Create(&application).Clauses(clause.Returning{}).Error; err != nil {
		return models.ApplicationCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return application, nil
}

func (a ApplicationGatewayImpl) DoesExistApplication(userId uint, nomination string) (bool, error) {
	if err := a.postgresClient.Db.Where("author_id = ? AND nomination = ?", userId, nomination).
		Take(&models.ApplicationCore{}).Error; err != nil {
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

func (a ApplicationGatewayImpl) GetApplicationsByAuthorId(id uint, offset, limit int) (applications []models.ApplicationCore, countRows uint, err error) {
	var count int64
	result := a.postgresClient.Db.Limit(limit).Offset(offset).Where("author_id = ?", id).
		Find(&applications)
	if result.Error != nil {
		return []models.ApplicationCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return applications, uint(count), result.Error
}
func (a ApplicationGatewayImpl) GetAllApplications(offset, limit int) (applications []models.ApplicationCore, countRows uint, err error) {
	var count int64
	result := a.postgresClient.Db.Limit(limit).Offset(offset).Find(&applications)
	if result.Error != nil {
		return []models.ApplicationCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return applications, uint(count), result.Error
}
