package models

import (
	"gorm.io/gorm"
	"time"
)

/*
	Таблица связи доступности мероприятия для региона
*/

type EventRegionRelCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	EventID   uint
	Event     EventCore `gorm:"foreignKey:EventID"`
	RegionID  uint
	Region    RegionCore `gorm:"foreignKey:RegionID"`
}
