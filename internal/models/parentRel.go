package models

import (
	"gorm.io/gorm"
	"time"
)

// TODO set gorm table name

type ParentRelCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ParentID  uint
	Parent    UserCore `gorm:"constraint:onUpdate:CASCADE;"`
	ChildID   uint
	Child     UserCore `gorm:"constraint:onUpdate:CASCADE;"`
}
