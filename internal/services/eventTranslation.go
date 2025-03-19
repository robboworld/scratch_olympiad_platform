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
	eventUserGateway        gateways.EventUserRelGateway
}

func (e EventTranslationServiceImpl) CreateEventTranslation(
	eventTranslation models.EventTranslationCore,
	clientId uint, clientRole models.Role,
) (
	newEventTranslation models.EventTranslationCore, err error,
) {
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		// Доступ только у SuperAdmin, Admin, Organizer
		clientEventRoles, err := e.eventUserGateway.GetEventRoles(eventTranslation.EventID, clientId)
		if err != nil {
			return models.EventTranslationCore{}, err
		}
		allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
		if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
			return models.EventTranslationCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
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
		// Доступ только у SuperAdmin, Admin, Organizer
		clientEventRoles, err := e.eventUserGateway.GetEventRoles(eventTranslationCore.EventID, clientId)
		if err != nil {
			return models.EventTranslationCore{}, err
		}
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
		// Доступ только у SuperAdmin, Admin, Organizer
		clientEventRoles, err := e.eventUserGateway.GetEventRoles(eventTranslation.EventID, clientId)
		if err != nil {
			return err
		}
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
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		// Доступ только у SuperAdmin, Admin, Organizer
		clientEventRoles, err := e.eventUserGateway.GetEventRoles(eventId, clientId)
		if err != nil {
			return models.EventTranslationCore{}, err
		}
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
