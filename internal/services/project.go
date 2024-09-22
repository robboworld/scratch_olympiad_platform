package services

import (
	"errors"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type ProjectService interface {
	UpdateProject(project models.ProjectCore) (updatedProject models.ProjectCore, err error)
	GetProjectById(id, clientId uint, clientRole models.Role) (project models.ProjectCore, err error)
}

type ProjectServiceImpl struct {
	projectGateway gateways.ProjectGateway
}

func (p ProjectServiceImpl) UpdateProject(project models.ProjectCore) (updatedProject models.ProjectCore, err error) {
	//TODO check is author?
	return p.projectGateway.UpdateProject(project)
}

func (p ProjectServiceImpl) GetProjectById(id, clientId uint, clientRole models.Role) (project models.ProjectCore, err error) {
	project, err = p.projectGateway.GetProjectById(id)
	if err != nil {
		return models.ProjectCore{}, err
	}
	// если проект забанен по решению супер админа, то доступ имеет только супер админ
	if project.IsBanned && clientRole.String() == models.RoleSuperAdmin.String() {
		return project, nil
	} else if project.IsBanned {
		return models.ProjectCore{}, errors.New(consts.ErrProjectPageIsBanned)
	}
	if project.IsShared || clientRole.String() == models.RoleSuperAdmin.String() {
		return project, nil
	} else {
		if project.AuthorID != clientId {
			return models.ProjectCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}
	return project, nil
}
