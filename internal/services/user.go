package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type UserService interface {
	CreateUser(user models.UserCore, clientRole models.Role) (newUser models.UserCore, err error)
	DeleteUser(id uint) error
	UpdateUser(user models.UserCore, clientRole models.Role) (updatedUser models.UserCore, err error)
	GetUserById(id uint, clientRole models.Role) (models.UserCore, error)
	GetAllUsers(page, pageSize *int, isActive bool, role []models.Role, clientRole models.Role) (users []models.UserCore, countRows uint, err error)
	SetIsActive(id uint, isActive bool) error
}

type UserServiceImpl struct {
	userGateway gateways.UserGateway
}

func (u UserServiceImpl) SetIsActive(id uint, isActive bool) error {
	return u.userGateway.SetIsActive(id, isActive)
}

// TODO use gorm hooks beforeCreate for example for check client role

func (u UserServiceImpl) CreateUser(user models.UserCore, clientRole models.Role) (newUser models.UserCore, err error) {
	// TODO сразу активен? надо ли высылать код
	// checking the client role for the possibility of creating a user
	if (clientRole == models.RoleUnitAdmin && user.Role.String() == models.RoleUnitAdmin.String()) ||
		user.Role.String() == models.RoleSuperAdmin.String() {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	exist, err := u.userGateway.DoesExistEmail(0, user.Email)
	if err != nil {
		return models.UserCore{}, err
	}
	if exist {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrEmailAlreadyInUse,
		}
	}
	if len(user.Password) < 6 {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrShortPassword,
		}
	}
	passwordHash := utils.HashPassword(user.Password)
	user.Password = passwordHash
	return u.userGateway.CreateUser(user)
}

func (u UserServiceImpl) DeleteUser(id uint) error {
	// TODO maybe check userById and role
	// TODO delete projects
	// delete all rels (projectPage, project, etc.)
	return u.userGateway.DeleteUser(id)
}

func (u UserServiceImpl) UpdateUser(user models.UserCore, clientRole models.Role) (updatedUser models.UserCore, err error) {
	// TODO check client role
	// какие роли кого могут обновлять
	exist, err := u.userGateway.DoesExistEmail(user.ID, user.Email)
	if err != nil {
		return models.UserCore{}, err
	}
	// checking the client role for the possibility of updating a user
	switch clientRole {
	case models.RoleUnitAdmin:
		if user.Role.String() == models.RoleSuperAdmin.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	case models.RoleTeacher:
		if user.Role.String() != models.RoleTeacher.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	case models.RoleParent:
		if user.Role.String() != models.RoleParent.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	case models.RoleStudent:
		if user.Role.String() != models.RoleStudent.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}
	if exist {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrEmailAlreadyInUse,
		}
	}
	return u.userGateway.UpdateUser(user)
}

func (u UserServiceImpl) GetUserById(id uint, clientRole models.Role) (models.UserCore, error) {
	user, err := u.userGateway.GetUserById(id)
	if err != nil {
		return models.UserCore{}, err
	}
	// checking the client role for the possibility of getting a user
	switch clientRole {
	case models.RoleParent:
		if user.Role.String() != models.RoleStudent.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	case models.RoleTeacher:
		if user.Role.String() == models.RoleUnitAdmin.String() ||
			user.Role.String() == models.RoleSuperAdmin.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	case models.RoleUnitAdmin:
		if user.Role.String() == models.RoleSuperAdmin.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}
	return user, nil
}

func (u UserServiceImpl) GetAllUsers(
	page, pageSize *int,
	isActive bool,
	roles []models.Role,
	clientRole models.Role,
) (users []models.UserCore, countRows uint, err error) {
	// checking the client role for the possibility of getting a users
	switch clientRole {
	case models.RoleUnitAdmin:
		for _, role := range roles {
			if role.String() == models.RoleSuperAdmin.String() || role.String() == models.RoleUnitAdmin.String() {
				return []models.UserCore{}, 0, utils.ResponseError{
					Code:    http.StatusForbidden,
					Message: consts.ErrAccessDenied,
				}
			}
		}
	}
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return u.userGateway.GetAllUsers(offset, limit, isActive, roles)
}
