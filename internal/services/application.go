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
	GetApplicationsByAuthorId(id, clientId uint, clientRole models.Role, page, pageSize *int) (applications []models.ApplicationCore, countRows uint, err error)
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

	subject := "Your submitted Scratch Olympiad application"
	body := "<p>Application details:</p>" +
		"<p>Nomination: " + application.Nomination + "</p>"

	if application.AlgorithmicTaskLink != "" {
		body += "<p>Algorithmic task link: " + application.AlgorithmicTaskLink + "</p>"
	}
	if application.AlgorithmicTaskFile != "" {
		body += "<p>Algorithmic task file: " + application.AlgorithmicTaskFile + "</p>"
	}
	if application.CreativeTaskLink != "" {
		body += "<p>Creative task link: " + application.CreativeTaskLink + "</p>"
	}
	if application.CreativeTaskFile != "" {
		body += "<p>Creative task file: " + application.CreativeTaskFile + "</p>"
	}
	if application.EngineeringTaskFile != "" {
		body += "<p>Engineering task file: " + application.EngineeringTaskFile + "</p>"
	}
	if application.EngineeringTaskCloudLink != "" {
		body += "<p>Engineering task cloud link: " + application.EngineeringTaskCloudLink + "</p>"
	}
	if application.EngineeringTaskVideo != "" {
		body += "<p>Engineering task video: " + application.EngineeringTaskVideo + "</p>"
	}
	if application.EngineeringTaskVideoCloudLink != "" {
		body += "<p>Engineering task video cloud link: " + application.EngineeringTaskVideoCloudLink + "</p>"
	}
	if application.Note != "" {
		body += "<p>Note: " + application.Note + "</p>"
	}

	body += "<br><p>Organizing committee of the International Scratch Creative Programming Olympiad</p>" +
		"<p><a href='mailto:scratch@creativeprogramming.org'>scratch@creativeprogramming.org</a></p>" +
		"<p><a href='https://creativeprogramming.org'>creativeprogramming.org</a></p>"

	if err = utils.SendEmail(subject, user.Email, body); err != nil {
		return models.ApplicationCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return a.applicationGateway.CreateApplication(application)
}

func (a ApplicationServiceImpl) GetApplicationsByAuthorId(id, clientId uint, clientRole models.Role, page, pageSize *int) (applications []models.ApplicationCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	if clientRole.String() != models.RoleSuperAdmin.String() {
		if id != clientId {
			return nil, 0, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
		return a.applicationGateway.GetApplicationsByAuthorId(id, offset, limit)
	}
	return a.applicationGateway.GetApplicationsByAuthorId(id, offset, limit)
}
