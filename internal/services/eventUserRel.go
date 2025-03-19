package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type EventUserService interface {
	CreateRel(eventId, userId uint, eventRole models.EventRole, clientId uint, clientRole models.Role) (models.EventUserRelCore, error)
	DeleteRel(eventId, userId uint, eventRole models.EventRole, clientId uint, clientRole models.Role) error

	GetEventRoles(eventId, userId uint) ([]models.EventRole, error)
}

type EventUserServiceImpl struct {
	userGateway         gateways.UserGateway
	eventUserGateway    gateways.EventUserRelGateway
	eventCountryGateway gateways.EventCountryRelGateway
	eventRegionGateway  gateways.EventRegionRelGateway
}

func (e EventUserServiceImpl) CreateRel(
	eventId, userId uint,
	eventRole models.EventRole,
	clientId uint, clientRole models.Role,
) (models.EventUserRelCore, error) {
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		switch eventRole {
		case models.EventRoleOrganizer:
			// Доступ только у SuperAdmin или Admin
			return models.EventUserRelCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		case models.EventRoleModerator, models.EventRoleExpert:
			// Доступ только у SuperAdmin, Admin или Organizer
			clientEventRoles, err := e.eventUserGateway.GetEventRoles(eventId, clientId)
			if err != nil {
				return models.EventUserRelCore{}, err
			}
			allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
			if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
				return models.EventUserRelCore{}, utils.ResponseError{
					Code:    http.StatusForbidden,
					Message: consts.ErrAccessDenied,
				}
			}
		}
	}

	// Нельзя назначить себя
	if userId == clientId {
		return models.EventUserRelCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrSelfAssignmentNotAllowed,
		}
	}

	user, err := e.userGateway.GetUserById(userId)
	if err != nil {
		return models.EventUserRelCore{}, err
	}

	// Нельзя назначить не активного пользователя
	if !user.IsActive {
		return models.EventUserRelCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrUserIsNotActive,
		}
	}

	// Нельзя назначить SuperAdmin и Admin
	if user.Role == models.RoleSuperAdmin || user.Role == models.RoleAdmin {
		return models.EventUserRelCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	// проверка доступности мероприятия для user по стране
	exist, err := e.eventCountryGateway.DoesExistRel(models.EventCountryRelCore{
		EventID:   eventId,
		CountryID: user.CountryID,
	})
	if err != nil {
		return models.EventUserRelCore{}, err
	}
	if !exist {
		return models.EventUserRelCore{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrEventNotAccessible,
		}
	}

	// проверка доступности мероприятия для user по региону
	if user.Country.HasRegions {
		if user.RegionID == nil {
			return models.EventUserRelCore{}, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrEventNotAccessible,
			}
		}
		exist, err = e.eventRegionGateway.DoesExistRel(models.EventRegionRelCore{
			EventID:  eventId,
			RegionID: *user.RegionID,
		})
		if err != nil {
			return models.EventUserRelCore{}, err
		}
		if !exist {
			return models.EventUserRelCore{}, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrEventNotAccessible,
			}
		}
	}

	return e.eventUserGateway.CreateRel(models.EventUserRelCore{
		EventID:   eventId,
		UserID:    userId,
		EventRole: eventRole,
	})
}

func (e EventUserServiceImpl) DeleteRel(
	eventId, userId uint,
	eventRole models.EventRole,
	clientId uint, clientRole models.Role,
) error {
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		switch eventRole {
		case models.EventRoleOrganizer:
			// Доступ только у SuperAdmin или Admin
			return utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		case models.EventRoleModerator, models.EventRoleExpert:
			clientEventRoles, err := e.eventUserGateway.GetEventRoles(eventId, clientId)
			if err != nil {
				return err
			}
			// Доступ только у SuperAdmin, Admin или Organizer
			allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
			if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
				return utils.ResponseError{
					Code:    http.StatusForbidden,
					Message: consts.ErrAccessDenied,
				}
			}
		}
	}
	return e.eventUserGateway.DeleteRel(models.EventUserRelCore{
		EventID:   eventId,
		UserID:    userId,
		EventRole: eventRole,
	})
}

func (e EventUserServiceImpl) GetEventRoles(eventId, userId uint) ([]models.EventRole, error) {
	return e.eventUserGateway.GetEventRoles(eventId, userId)
}
