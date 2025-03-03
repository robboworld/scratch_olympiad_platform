package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type EventTranslationCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`

	EventID uint
	Event   EventCore `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;'"`
}

func (e *EventTranslationHTTP) FromCore(eventTranslation EventTranslationCore) {
	e.ID = strconv.Itoa(int(eventTranslation.ID))
	e.Name = eventTranslation.Name
	e.Description = eventTranslation.Description
	e.EventID = strconv.Itoa(int(eventTranslation.EventID))
}
