package models

import (
	"gorm.io/gorm"
	"time"
)

type SolutionCore struct {
	ID         uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Name       string         `gorm:"unique;size:255"`
	AccessLink string
}
