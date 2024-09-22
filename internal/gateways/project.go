package gateways

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"gorm.io/gorm/clause"
	"net/http"
)

type ProjectGateway interface {
	CreateProject(project models.ProjectCore) (models.ProjectCore, error)
	DeleteProject(id, clientId uint) error
	UpdateProject(project models.ProjectCore) (updatedProject models.ProjectCore, err error)
	GetProjectById(id uint) (project models.ProjectCore, err error)
	GetProjectsByAuthorId(id uint) (projects []models.ProjectCore, err error)
}

type ProjectGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (p ProjectGatewayImpl) GetProjectsByAuthorId(id uint) (projects []models.ProjectCore, err error) {
	if err := p.postgresClient.Db.Where("author_id = ?", id).Find(&projects).Error; err != nil {
		return []models.ProjectCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return projects, nil
}

func (p ProjectGatewayImpl) CreateProject(project models.ProjectCore) (models.ProjectCore, error) {
	result := p.postgresClient.Db.Create(&project).Clauses(clause.Returning{})
	return project, result.Error
}

func (p ProjectGatewayImpl) DeleteProject(id, clientId uint) error {
	result := p.postgresClient.Db.Unscoped().Where("author_id = ?", clientId).Delete(&models.ProjectCore{ID: id})
	var countRows int
	if err := result.Row().Scan(&countRows); err != nil {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	if countRows == 0 && result.Error == nil {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrNotFoundInDB,
		}
	}
	return result.Error
}

func (p ProjectGatewayImpl) UpdateProject(project models.ProjectCore) (updatedProject models.ProjectCore, err error) {
	result := p.postgresClient.Db.Where(&models.ProjectCore{ID: project.ID}).Updates(project)
	return project, result.Error
}

func (p ProjectGatewayImpl) GetProjectById(id uint) (project models.ProjectCore, err error) {
	if err := p.postgresClient.Db.First(&project, id).Error; err != nil {
		return models.ProjectCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return project, nil
}
