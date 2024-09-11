package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UserCore struct {
	ID                  uint `gorm:"primaryKey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	Email               string         `gorm:"not null;"`
	Password            string         `gorm:"not null;"`
	Role                Role           `gorm:"not null;"`
	FullName            string         `gorm:"not null;"`
	FullNameNative      string         `gorm:"not null;"`
	Country             string         `gorm:"not null;"`
	City                string         `gorm:"not null;"`
	Birthdate           time.Time      `gorm:"not null;"`
	IsActive            bool           `gorm:"not null;default:false;type:boolean;column:is_active"`
	ActivationLink      string
	VerificationCode    string
	VerificationCodeTtl time.Time
}

func (u *UserHTTP) ToCore() UserCore {
	id, _ := strconv.ParseUint(u.ID, 10, 64)
	birthDate, _ := time.Parse(time.DateOnly, u.Birthdate)
	return UserCore{
		ID:             uint(id),
		Email:          u.Email,
		Password:       u.Password,
		Role:           u.Role,
		FullName:       u.FullName,
		FullNameNative: u.FullNameNative,
		Country:        u.Country,
		City:           u.City,
		Birthdate:      birthDate,
		IsActive:       u.IsActive,
		ActivationLink: u.ActivationLink,
	}
}

func (u *UserHTTP) FromCore(userCore UserCore) {
	u.ID = strconv.Itoa(int(userCore.ID))
	u.CreatedAt = userCore.CreatedAt.Format(time.DateTime)
	u.UpdatedAt = userCore.UpdatedAt.Format(time.DateTime)
	u.Email = userCore.Email
	u.FullName = userCore.FullName
	u.FullNameNative = userCore.FullNameNative
	u.Country = userCore.Country
	u.City = userCore.City
	u.Birthdate = userCore.Birthdate.Format(time.DateOnly)
	u.IsActive = userCore.IsActive
	u.Role = userCore.Role
}

func FromUsersCore(usersCore []UserCore) (usersHttp []*UserHTTP) {
	for _, userCore := range usersCore {
		var tmpUserHttp UserHTTP
		tmpUserHttp.FromCore(userCore)
		usersHttp = append(usersHttp, &tmpUserHttp)
	}
	return
}
