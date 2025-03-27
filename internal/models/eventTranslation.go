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
	EventID   uint
	Event     EventCore

	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
}

func (e *EventTranslationHTTP) FromCore(translation EventTranslationCore) {
	e.ID = strconv.Itoa(int(translation.ID))
	e.CreatedAt = translation.CreatedAt.Format(time.DateTime)
	e.UpdatedAt = translation.UpdatedAt.Format(time.DateTime)
	e.Name = translation.Name
	e.Description = translation.Description
	e.EventID = strconv.Itoa(int(translation.EventID))
}
