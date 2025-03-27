package models

import (
	"gorm.io/gorm"
	"time"
)

/*
	Таблица связи доступа пользователя по роли к мероприятию
*/

type EventUserRelCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	EventID   uint
	Event     EventCore `gorm:"foreignKey:EventID"`
	UserID    uint
	User      UserCore  `gorm:"foreignKey:UserID"`
	EventRole EventRole `gorm:"not null"`
}
