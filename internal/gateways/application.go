package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
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
	result := u.postgresClient.Db.Create(&application).Clauses(clause.Returning{})
	if result.Error != nil {
		return models.ApplicationCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	return application, nil
}
