package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/api"
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
	applicationAPI     api.ApplicationAPI
}

func (a ApplicationServiceImpl) CreateApplication(application models.ApplicationCore) (models.ApplicationCore, error) {
	user, err := a.userGateway.GetUserById(application.AuthorID)
	if err != nil {
		return models.ApplicationCore{}, err
	}
	nomination, err := a.nominationGateway.GetNominationByName(application.Nomination)
	if err != nil {
		return models.ApplicationCore{}, err
	}

	exist, err := a.applicationGateway.DoesExistApplication(application.AuthorID, application.Nomination)
	if err != nil {
		return models.ApplicationCore{}, err
	}
	if exist {
		return models.ApplicationCore{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrApplicationAlreadySubmitted,
		}
	}

	userAge := uint(utils.CalculateUserAge(user.Birthdate))
	if userAge < nomination.MinAge {
		return models.ApplicationCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrDoesNotMatchAgeCategory,
		}
	}

	application.Author = user
	err = a.applicationAPI.CreateApplication(application)
	if err != nil {
		return models.ApplicationCore{}, err
	}

	return a.applicationGateway.CreateApplication(application)
}
