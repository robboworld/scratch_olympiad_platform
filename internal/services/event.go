package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
)

type EventService interface {
	CreateEvent(newEvent models.EventCore) (models.EventCore, error)
	UpdateEvent(event models.EventCore) (updatedEvent models.EventCore, err error)
	GetEventById(id uint) (event models.EventCore, err error)
	GetAllEvents(page, pageSize *int) (events []models.EventCore, countRows uint, err error)
}

type EventServiceImpl struct {
	eventGateway gateways.EventGateway
}

func (e EventServiceImpl) CreateEvent(newEvent models.EventCore) (event models.EventCore, err error) {
	return e.eventGateway.CreateEvent(newEvent)
}

func (e EventServiceImpl) UpdateEvent(event models.EventCore) (updatedEvent models.EventCore, err error) {
	return e.eventGateway.UpdateEvent(event)
}

func (e EventServiceImpl) GetEventById(id uint) (event models.EventCore, err error) {
	return e.eventGateway.GetEventById(id)
}

func (e EventServiceImpl) GetAllEvents(page, pageSize *int) (events []models.EventCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return e.eventGateway.GetAllEvents(offset, limit)
}
