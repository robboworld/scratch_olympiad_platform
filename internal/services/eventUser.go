package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type EventUserService interface {
	SetOrganizerForEvent(eventId, userId, clientId uint, clientRole models.Role) error
	SetModeratorForEvent(eventId, userId, clientId uint, clientRole models.Role) error
	SetExpertForEvent(eventId, userId, clientId uint, clientRole models.Role) error

	UnsetOrganizerForEvent(eventId, userId uint, clientRole models.Role) error
	UnsetModeratorForEvent(eventId, userId, clientId uint, clientRole models.Role) error
	UnsetExpertForEvent(eventId, userId, clientId uint, clientRole models.Role) error
}

type EventUserServiceImpl struct {
	userGateway         gateways.UserGateway
	eventGateway        gateways.EventGateway
	eventUserGateway    gateways.EventUserGateway
	eventCountryGateway gateways.EventCountryGateway
	eventRegionGateway  gateways.EventRegionGateway
}

func (e EventUserServiceImpl) SetOrganizerForEvent(eventId, userId, clientId uint, clientRole models.Role) error {
	// Нельзя назначить себя организатором
	if userId == clientId {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrSelfAssignmentNotAllowed,
		}
	}

	user, err := e.userGateway.GetUserById(userId)
	if err != nil {
		return err
	}

	// Нельзя назначить не активного пользователя
	if !user.IsActive {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrUserIsNotActive,
		}
	}

	// Нельзя назначить SuperAdmin и Admin - организатором
	if user.Role.String() == models.RoleSuperAdmin.String() || user.Role.String() == models.RoleAdmin.String() {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	// Нет доступа для User
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	_, err = e.eventGateway.GetEventById(eventId)
	if err != nil {
		return err
	}

	// проверка доступности мероприятия для user по стране
	exist, err := e.eventCountryGateway.DoesExistEventCountry(eventId, user.CountryID)
	if err != nil {
		return err
	}
	if !exist {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrEventNotAccessible,
		}
	}

	// проверка доступности мероприятия для user по региону
	if user.Country.HasRegions {
		exist, err = e.eventRegionGateway.DoesExistEventRegion(eventId, *user.RegionID)
		if err != nil {
			return err
		}
		if !exist {
			return utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrEventNotAccessible,
			}
		}
	}

	return e.eventUserGateway.SetUserForEvent(eventId, userId, models.EventRoleOrganizer)
}

func (e EventUserServiceImpl) SetModeratorForEvent(eventId, userId, clientId uint, clientRole models.Role) error {
	// Нельзя назначить себя модератором
	if userId == clientId {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrSelfAssignmentNotAllowed,
		}
	}

	user, err := e.userGateway.GetUserById(userId)
	if err != nil {
		return err
	}

	// Нельзя назначить не активного пользователя
	if !user.IsActive {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrUserIsNotActive,
		}
	}

	// Нельзя назначить SuperAdmin и Admin - модератором
	if user.Role.String() == models.RoleSuperAdmin.String() || user.Role.String() == models.RoleAdmin.String() {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		clientEventRoles, err := e.eventUserGateway.GetUserRolesForEvent(eventId, clientId)
		if err != nil {
			return err
		}
		// Если у пользователя нет роли Organizer, доступа нет
		allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
		if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
			return utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}

	_, err = e.eventGateway.GetEventById(eventId)
	if err != nil {
		return err
	}
	// проверка доступности мероприятия для user по стране
	exist, err := e.eventCountryGateway.DoesExistEventCountry(eventId, user.CountryID)
	if err != nil {
		return err
	}
	if !exist {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrEventNotAccessible,
		}
	}
	// проверка доступности мероприятия для user по региону
	if user.Country.HasRegions {
		exist, err = e.eventRegionGateway.DoesExistEventRegion(eventId, *user.RegionID)
		if err != nil {
			return err
		}
		if !exist {
			return utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrEventNotAccessible,
			}
		}
	}
	return e.eventUserGateway.SetUserForEvent(eventId, userId, models.EventRoleModerator)
}

func (e EventUserServiceImpl) SetExpertForEvent(eventId, userId, clientId uint, clientRole models.Role) error {
	// Нельзя назначить себя экспертом
	if userId == clientId {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrSelfAssignmentNotAllowed,
		}
	}

	user, err := e.userGateway.GetUserById(userId)
	if err != nil {
		return err
	}

	// Нельзя назначить не активного пользователя
	if !user.IsActive {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrUserIsNotActive,
		}
	}

	// Нельзя назначить SuperAdmin и Admin - экспертом
	if user.Role.String() == models.RoleSuperAdmin.String() || user.Role.String() == models.RoleAdmin.String() {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		clientEventRoles, err := e.eventUserGateway.GetUserRolesForEvent(eventId, clientId)
		if err != nil {
			return err
		}
		// Если у пользователя нет роли Organizer, доступа нет
		allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
		if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
			return utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}
	_, err = e.eventGateway.GetEventById(eventId)
	if err != nil {
		return err
	}
	// проверка доступности мероприятия для user по стране
	exist, err := e.eventCountryGateway.DoesExistEventCountry(eventId, user.CountryID)
	if err != nil {
		return err
	}
	if !exist {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrEventNotAccessible,
		}
	}
	// проверка доступности мероприятия для user по региону
	if user.Country.HasRegions {
		exist, err = e.eventRegionGateway.DoesExistEventRegion(eventId, *user.RegionID)
		if err != nil {
			return err
		}
		if !exist {
			return utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrEventNotAccessible,
			}
		}
	}
	return e.eventUserGateway.SetUserForEvent(eventId, userId, models.EventRoleExpert)
}

func (e EventUserServiceImpl) UnsetOrganizerForEvent(eventId, userId uint, clientRole models.Role) error {
	// Нет доступа для User
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}
	return e.eventUserGateway.UnsetUserForEvent(eventId, userId, models.EventRoleOrganizer)
}

func (e EventUserServiceImpl) UnsetModeratorForEvent(eventId, userId, clientId uint, clientRole models.Role) error {
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		clientEventRoles, err := e.eventUserGateway.GetUserRolesForEvent(eventId, clientId)
		if err != nil {
			return err
		}
		// Если у пользователя нет роли Organizer, доступа нет
		allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
		if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
			return utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}

	return e.eventUserGateway.UnsetUserForEvent(eventId, userId, models.EventRoleModerator)
}

func (e EventUserServiceImpl) UnsetExpertForEvent(eventId, userId, clientId uint, clientRole models.Role) error {
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		clientEventRoles, err := e.eventUserGateway.GetUserRolesForEvent(eventId, clientId)
		if err != nil {
			return err
		}
		// Если у пользователя нет роли Organizer, доступа нет
		allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
		if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
			return utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}

	return e.eventUserGateway.UnsetUserForEvent(eventId, userId, models.EventRoleExpert)
}
