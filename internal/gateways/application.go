package gateways

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"gorm.io/gorm/clause"
	"net/http"
)

type ApplicationGateway interface {
	CreateApplication(application models.ApplicationCore) (newApplication models.ApplicationCore, err error)
}

type ApplicationGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (u ApplicationGatewayImpl) CreateApplication(application models.ApplicationCore) (newApplication models.ApplicationCore, err error) {
	if err = u.postgresClient.Db.Create(&application).Clauses(clause.Returning{}).Error; err != nil {
		return models.ApplicationCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return application, nil
}
