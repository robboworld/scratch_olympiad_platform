package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type EventTranslationService interface {
	CreateEventTranslation(eventTranslation models.EventTranslationCore, clientId uint, clientRole models.Role) (newEventTranslation models.EventTranslationCore, err error)
	UpdateEventTranslation(eventTranslation models.EventTranslationCore, clientId uint, clientRole models.Role) (updatedEventTranslation models.EventTranslationCore, err error)
	DeleteEventTranslation(id uint, clientId uint, clientRole models.Role) error

	// GetEventTranslationByEventId TODO: Когда будет несколько переводов, переделать под получения списка переводов
	GetEventTranslationByEventId(eventId uint, clientId uint, clientRole models.Role) (eventTranslation models.EventTranslationCore, err error)
}

type EventTranslationServiceImpl struct {
	eventGateway            gateways.EventGateway
	eventTranslationGateway gateways.EventTranslationGateway
	eventUserGateway        gateways.EventUserGateway
}

func (e EventTranslationServiceImpl) CreateEventTranslation(
	eventTranslation models.EventTranslationCore,
	clientId uint, clientRole models.Role,
) (
	newEventTranslation models.EventTranslationCore, err error,
) {
	_, err = e.eventGateway.GetEventById(eventTranslation.EventID)
	if err != nil {
		return models.EventTranslationCore{}, err
	}

	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		clientEventRoles, err := e.eventUserGateway.GetUserRolesForEvent(eventTranslation.EventID, clientId)
		if err != nil {
			return models.EventTranslationCore{}, err
		}
		// Если у пользователя нет роли Organizer, доступа нет
		allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
		if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
			return models.EventTranslationCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}

	exists, err := e.eventTranslationGateway.DoesExistEventTranslation(eventTranslation.EventID)
	if err != nil {
		return models.EventTranslationCore{}, err
	}
	if exists {
		return models.EventTranslationCore{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrTranslationForEventAlreadyExist,
		}
	}

	return e.eventTranslationGateway.CreateEventTranslation(eventTranslation)
}

func (e EventTranslationServiceImpl) UpdateEventTranslation(
	eventTranslation models.EventTranslationCore,
	clientId uint, clientRole models.Role,
) (
	updatedEventTranslation models.EventTranslationCore, err error,
) {
	eventTranslationCore, err := e.eventTranslationGateway.GetEventTranslationById(eventTranslation.ID)
	if err != nil {
		return models.EventTranslationCore{}, err
	}

	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		clientEventRoles, err := e.eventUserGateway.GetUserRolesForEvent(eventTranslationCore.EventID, clientId)
		if err != nil {
			return models.EventTranslationCore{}, err
		}
		// Если у пользователя нет роли Organizer, доступа нет
		allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
		if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
			return models.EventTranslationCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}

	return e.eventTranslationGateway.UpdateEventTranslation(eventTranslation)
}

func (e EventTranslationServiceImpl) DeleteEventTranslation(id uint, clientId uint, clientRole models.Role) error {
	eventTranslation, err := e.eventTranslationGateway.GetEventTranslationById(id)
	if err != nil {
		return err
	}

	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		clientEventRoles, err := e.eventUserGateway.GetUserRolesForEvent(eventTranslation.EventID, clientId)
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
	return e.eventTranslationGateway.DeleteEventTranslation(id)
}

func (e EventTranslationServiceImpl) GetEventTranslationByEventId(eventId uint, clientId uint, clientRole models.Role) (
	eventTranslation models.EventTranslationCore, err error,
) {
	_, err = e.eventGateway.GetEventById(eventId)
	if err != nil {
		return models.EventTranslationCore{}, err
	}

	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		clientEventRoles, err := e.eventUserGateway.GetUserRolesForEvent(eventId, clientId)
		if err != nil {
			return models.EventTranslationCore{}, err
		}
		// Если у пользователя нет роли Organizer, доступа нет
		allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
		if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
			return models.EventTranslationCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}

	return e.eventTranslationGateway.GetEventTranslationByEventId(eventId)
}
