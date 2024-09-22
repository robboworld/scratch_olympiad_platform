package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type ProjectPageService interface {
	CreateProjectPage(authorId uint) (newProjectPage models.ProjectPageCore, err error)
	DeleteProjectPage(id, clientId uint) error
	GetAllProjectPages(page, pageSize *int, userId uint, clientRole models.Role) (projectPages []models.ProjectPageCore, countRows uint, err error)
	UpdateProjectPage(projectPage models.ProjectPageCore, clientId uint) (models.ProjectPageCore, error)
	GetProjectPageById(id, clientId uint, clientRole models.Role) (projectPage models.ProjectPageCore, err error)
	GetProjectsPageByAuthorId(id uint, page, pageSize *int) (projectPages []models.ProjectPageCore, countRows uint, err error)
	SetIsBanned(id uint, isBanned bool) error
}

type ProjectPageServiceImpl struct {
	projectGateway     gateways.ProjectGateway
	projectPageGateway gateways.ProjectPageGateway
}

func (p ProjectPageServiceImpl) SetIsBanned(id uint, isBanned bool) error {
	return p.projectPageGateway.SetIsBanned(id, isBanned)
}

func (p ProjectPageServiceImpl) GetAllProjectPages(page, pageSize *int, userId uint, clientRole models.Role) (projectPages []models.ProjectPageCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	if clientRole.String() != models.RoleSuperAdmin.String() {
		return p.projectPageGateway.GetProjectPagesByAuthorId(userId, offset, limit)
	}
	return p.projectPageGateway.GetAllProjectPages(offset, limit)
}

func (p ProjectPageServiceImpl) GetProjectsPageByAuthorId(id uint, page, pageSize *int) (projectPages []models.ProjectPageCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return p.projectPageGateway.GetProjectPagesByAuthorId(id, offset, limit)
}

func (p ProjectPageServiceImpl) CreateProjectPage(authorId uint) (newProjectPage models.ProjectPageCore, err error) {
	return p.projectPageGateway.CreateProjectPage(
		models.ProjectPageCore{
			AuthorID:    authorId,
			Title:       "Untitled",
			Instruction: "",
			Notes:       "",
			IsShared:    false,
		},
		models.ProjectCore{
			AuthorID: authorId,
			Json:     consts.EmptyProjectJson,
		})
}

func (p ProjectPageServiceImpl) DeleteProjectPage(id, clientId uint) error {
	return p.projectPageGateway.DeleteProjectPage(id, clientId)
}

func (p ProjectPageServiceImpl) UpdateProjectPage(projectPage models.ProjectPageCore, clientId uint) (models.ProjectPageCore, error) {
	getProjectPage, err := p.projectPageGateway.GetProjectPageById(projectPage.ID)
	if err != nil {
		return models.ProjectPageCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	if clientId != getProjectPage.AuthorID {
		return models.ProjectPageCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}
	return p.projectPageGateway.UpdateProjectPage(projectPage)
}

func (p ProjectPageServiceImpl) GetProjectPageById(id, clientId uint, clientRole models.Role) (projectPage models.ProjectPageCore, err error) {
	projectPage, err = p.projectPageGateway.GetProjectPageById(id)
	if err != nil {
		return projectPage, err
	}
	// если проект забанен по решению супер админа, то доступ имеет только супер админ
	if projectPage.IsBanned && clientRole.String() == models.RoleSuperAdmin.String() {
		return projectPage, nil
	} else if projectPage.IsBanned {
		return models.ProjectPageCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrProjectPageIsBanned,
		}
	}
	// проверка доступа к проекту. супер админу всегда имеет доступ к проекту
	if projectPage.IsShared || clientRole.String() == models.RoleSuperAdmin.String() {
		return projectPage, nil
	} else {
		if projectPage.AuthorID != clientId {
			return models.ProjectPageCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}
	return projectPage, nil
}
