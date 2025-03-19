package models

import (
	"gorm.io/gorm"
	"time"
)

/*
	Таблица связи доступности мероприятия для страны
*/

type EventCountryRelCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	EventID   uint
	Event     EventCore `gorm:"foreignKey:EventID"`
	CountryID uint
	Country   CountryCore `gorm:"foreignKey:CountryID"`
}
