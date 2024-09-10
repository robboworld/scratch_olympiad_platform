package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type ApplicationService interface {
	CreateApplication(newApplication models.ApplicationCore) (models.ApplicationCore, error)
}

type ApplicationServiceImpl struct {
	applicationGateway gateways.ApplicationGateway
	nominationGateway  gateways.NominationGateway
	userGateway        gateways.UserGateway
}

func (a ApplicationServiceImpl) CreateApplication(application models.ApplicationCore) (models.ApplicationCore, error) {
	exist, err := a.nominationGateway.DoesExistName(0, application.Nomination)
	if err != nil {
		return models.ApplicationCore{}, err
	}
	if !exist {
		return models.ApplicationCore{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrNominationNotFoundInDB,
		}
	}
	client, err := a.userGateway.GetUserById(application.AuthorID)
	if err != nil {
		return models.ApplicationCore{}, err
	}

	nomination, err := a.nominationGateway.GetNominationByName(application.Nomination)
	if err != nil {
		return models.ApplicationCore{}, err
	}

	clientAge := uint(utils.CalculateAge(client.Birthdate))
	if clientAge < nomination.MinAge || clientAge > nomination.MaxAge {
		return models.ApplicationCore{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrDoesNotMatchAge,
		}
	}

	return a.applicationGateway.CreateApplication(application)
}
