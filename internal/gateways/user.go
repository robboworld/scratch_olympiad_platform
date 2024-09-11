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

type UserGateway interface {
	CreateUser(user models.UserCore) (newUser models.UserCore, err error)
	DeleteUser(id uint) (err error)
	UpdateUser(user models.UserCore) (updatedUser models.UserCore, err error)
	GetUserById(id uint) (user models.UserCore, err error)
	GetUserByActivationLink(link string) (user models.UserCore, err error)
	GetUserByEmail(email string) (user models.UserCore, err error)
	GetAllUsers(offset, limit int, isActive bool, role []models.Role) (users []models.UserCore, countRows uint, err error)
	DoesExistEmail(id uint, email string) (bool, error)
	SetIsActive(id uint, isActive bool) error
}

type UserGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (u UserGatewayImpl) GetUserByActivationLink(link string) (user models.UserCore, err error) {
	if err = u.postgresClient.Db.Where("activation_link = ?", link).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return user, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return user, nil
}

func (u UserGatewayImpl) GetUserByEmail(email string) (user models.UserCore, err error) {
	if err = u.postgresClient.Db.Where("email = ?", email).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrIncorrectPasswordOrEmail,
			}
		}
		return user, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return user, nil
}

func (u UserGatewayImpl) SetIsActive(id uint, isActive bool) (err error) {
	var updateStruct map[string]interface{}
	if isActive {
		updateStruct = map[string]interface{}{
			"is_active":       isActive,
			"activation_link": "",
		}
	} else {
		updateStruct = map[string]interface{}{
			"is_active": isActive,
		}
	}
	if err = u.postgresClient.Db.First(&models.UserCore{ID: id}).Updates(updateStruct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

func (u UserGatewayImpl) DoesExistEmail(id uint, email string) (bool, error) {
	if err := u.postgresClient.Db.Where("id != ? AND email = ?", id, email).
		Take(&models.UserCore{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return true, nil
}

func (u UserGatewayImpl) CreateUser(user models.UserCore) (newUser models.UserCore, err error) {
	if err = u.postgresClient.Db.Create(&user).Clauses(clause.Returning{}).Error; err != nil {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return user, nil
}

func (u UserGatewayImpl) DeleteUser(id uint) (err error) {
	if err = u.postgresClient.Db.Take(&models.UserCore{}, id).Delete(&models.UserCore{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return
}

func (u UserGatewayImpl) UpdateUser(user models.UserCore) (models.UserCore, error) {
	if err := u.postgresClient.Db.Model(&user).Clauses(clause.Returning{}).
		Take(&models.UserCore{}, user.ID).
		Updates(map[string]interface{}{
			"email":            user.Email,
			"full_name":        user.FullName,
			"full_name_native": user.FullNameNative,
			"country":          user.Country,
			"city":             user.City,
			"birthdate":        user.Birthdate,
		}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return user, nil
}

func (u UserGatewayImpl) GetUserById(id uint) (user models.UserCore, err error) {
	if err = u.postgresClient.Db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrNotFoundInDB,
			}
		}
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return user, nil
}

func (u UserGatewayImpl) GetAllUsers(
	offset, limit int,
	isActive bool,
	role []models.Role,
) (users []models.UserCore, countRows uint, err error) {
	var count int64
	if len(role) == 0 {
		role = append(role,
			models.RoleStudent,
			models.RoleParent,
			models.RoleTeacher,
			models.RoleUnitAdmin,
		)
	}
	result := u.postgresClient.Db.Limit(limit).Offset(offset).
		Where("is_active = ? AND (role) IN ?", isActive, role).Find(&users)
	if result.Error != nil {
		return []models.UserCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return users, uint(count), nil
}
