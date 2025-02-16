package models

import "time"

type AuthDataCore struct {
	ID uint `gorm:"primaryKey"`

	ActivationToken      string
	PasswordResetToken   string
	PasswordResetTokenAt time.Time

	UserID uint
	User   UserCore `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
