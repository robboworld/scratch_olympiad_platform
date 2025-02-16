package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type EventCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	StartDate   time.Time `gorm:"not null"`
	EndDate     time.Time `gorm:"not null"`
}

func (e *EventHTTP) FromCore(event EventCore) {
	e.ID = strconv.Itoa(int(event.ID))
	e.CreatedAt = event.CreatedAt.Format(time.DateTime)
	e.UpdatedAt = event.UpdatedAt.Format(time.DateTime)
	e.Name = event.Name
	e.StartDate = event.StartDate.Format(time.DateOnly)
	e.EndDate = event.EndDate.Format(time.DateOnly)
}

func (e *EventDetailsHTTP) FromCore(event EventCore) {
	e.ID = strconv.Itoa(int(event.ID))
	e.CreatedAt = event.CreatedAt.Format(time.DateTime)
	e.UpdatedAt = event.UpdatedAt.Format(time.DateTime)
	e.Name = event.Name
	e.Description = event.Description
	e.StartDate = event.StartDate.Format(time.DateOnly)
	e.EndDate = event.EndDate.Format(time.DateOnly)
}

func FromEventsCore(eventsCore []EventCore) (eventsHttp []*EventHTTP) {
	for _, eventCore := range eventsCore {
		var tmpEventHttp EventHTTP
		tmpEventHttp.FromCore(eventCore)
		eventsHttp = append(eventsHttp, &tmpEventHttp)
	}
	return
}
