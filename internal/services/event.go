package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"net/http"
)

type EventService interface {
	CreateEvent(newEvent models.EventCore) (event models.EventCore, err error)
	UpdateEvent(event models.EventCore, clientId uint, clientRole models.Role) (updatedEvent models.EventCore, err error)
	GetEventById(id, clientId uint, clientRole models.Role) (event models.EventCore, err error)
	GetOriginalEventById(id, clientId uint, clientRole models.Role) (event models.EventCore, err error)
	GetAllEvents(page, pageSize *int, clientId uint, clientRole models.Role) (events []models.EventCore, countRows uint, err error)
}

type EventServiceImpl struct {
	userGateway             gateways.UserGateway
	eventGateway            gateways.EventGateway
	eventUserGateway        gateways.EventUserGateway
	eventTranslationGateway gateways.EventTranslationGateway
	eventCountryGateway     gateways.EventCountryGateway
	eventRegionGateway      gateways.EventRegionGateway
}

func (e EventServiceImpl) CreateEvent(newEvent models.EventCore) (event models.EventCore, err error) {
	return e.eventGateway.CreateEvent(newEvent)
}

func (e EventServiceImpl) UpdateEvent(event models.EventCore, clientId uint, clientRole models.Role) (updatedEvent models.EventCore, err error) {
	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		clientEventRoles, err := e.eventUserGateway.GetUserRolesForEvent(event.ID, clientId)
		if err != nil {
			return models.EventCore{}, err
		}
		// Если у пользователя нет роли Organizer, доступа нет
		allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
		if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
			return models.EventCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}
	return e.eventGateway.UpdateEvent(event)
}

func (e EventServiceImpl) GetEventById(id, clientId uint, clientRole models.Role) (event models.EventCore, err error) {
	event, err = e.eventGateway.GetEventById(id)
	if err != nil {
		return models.EventCore{}, err
	}

	client, err := e.userGateway.GetUserById(clientId)
	if err != nil {
		return models.EventCore{}, err
	}

	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		// проверка доступности мероприятия для client по стране
		exist, err := e.eventCountryGateway.DoesExistEventCountry(event.ID, client.CountryID)
		if err != nil {
			return models.EventCore{}, err
		}
		if !exist {
			return models.EventCore{}, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrEventNotAccessible,
			}
		}
		// проверка доступности мероприятия для client по региону
		if client.Country.HasRegions {
			if client.RegionID == nil {
				return models.EventCore{}, utils.ResponseError{
					Code:    http.StatusBadRequest,
					Message: consts.ErrEventNotAccessible,
				}
			}
			exist, err = e.eventRegionGateway.DoesExistEventRegion(event.ID, *client.RegionID)
			if err != nil {
				return models.EventCore{}, err
			}
			if !exist {
				return models.EventCore{}, utils.ResponseError{
					Code:    http.StatusBadRequest,
					Message: consts.ErrEventNotAccessible,
				}
			}
		}
	}

	exists, err := e.eventTranslationGateway.DoesExistEventTranslation(id)
	if err != nil {
		return models.EventCore{}, err
	}
	if exists {
		eventTranslation, err := e.eventTranslationGateway.GetEventTranslationByEventId(id)
		if err != nil {
			return models.EventCore{}, err
		}
		event.Name = eventTranslation.Name
		event.Description = eventTranslation.Description
	}

	return event, nil
}

func (e EventServiceImpl) GetOriginalEventById(id, clientId uint, clientRole models.Role) (event models.EventCore, err error) {
	event, err = e.eventGateway.GetEventById(id)
	if err != nil {
		return models.EventCore{}, err
	}

	if clientRole != models.RoleSuperAdmin && clientRole != models.RoleAdmin {
		clientEventRoles, err := e.eventUserGateway.GetUserRolesForEvent(event.ID, clientId)
		if err != nil {
			return models.EventCore{}, err
		}
		// Если у пользователя нет роли Organizer, доступа нет
		allowedEventRoles := []models.EventRole{models.EventRoleOrganizer}
		if !utils.DoesHaveEventRole(clientEventRoles, allowedEventRoles) {
			return models.EventCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}

	return event, nil
}

func (e EventServiceImpl) GetAllEvents(
	page, pageSize *int,
	clientId uint,
	clientRole models.Role,
) (events []models.EventCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)

	// Получаем мероприятие в зависимости от роли пользователя
	if clientRole == models.RoleSuperAdmin || clientRole == models.RoleAdmin {
		events, countRows, err = e.eventGateway.GetAllEvents(offset, limit)
	} else {
		client, err := e.userGateway.GetUserById(clientId)
		if err != nil {
			return nil, 0, err
		}
		// Получаем мероприятия доступные для страны и региона
		if client.Country.HasRegions {
			events, countRows, err = e.eventGateway.GetEventsByCountryIdAndRegionId(client.CountryID, *client.RegionID, offset, limit)
		} else {
			events, countRows, err = e.eventGateway.GetEventsByCountryId(client.CountryID, offset, limit)
		}
	}
	if err != nil {
		return nil, 0, err
	}

	for i := range events {
		exists, err := e.eventTranslationGateway.DoesExistEventTranslation(events[i].ID)
		if err != nil {
			return nil, 0, err
		}
		if exists {
			eventTranslation, err := e.eventTranslationGateway.GetEventTranslationByEventId(events[i].ID)
			if err != nil {
				return nil, 0, err
			}
			events[i].Name = eventTranslation.Name
			events[i].Description = eventTranslation.Description
		}
	}

	return events, countRows, nil
}
