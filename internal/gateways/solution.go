package gateways

import (
	"errors"
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/db"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type SolutionGateway interface {
	CreateSolution(solution models.SolutionCore) (newSolution models.SolutionCore, err error)
	GetSolutionByName(name string) (solution models.SolutionCore, err error)
}

type SolutionGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (s SolutionGatewayImpl) CreateSolution(solution models.SolutionCore) (newSolution models.SolutionCore, err error) {
	if err = s.postgresClient.Db.Create(&solution).Clauses(clause.Returning{}).Error; err != nil {
		return models.SolutionCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return solution, nil
}

func (s SolutionGatewayImpl) GetSolutionByName(name string) (solution models.SolutionCore, err error) {
	if err = s.postgresClient.Db.Where("name = ?", name).Take(&solution).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return solution, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return solution, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return solution, nil
}
