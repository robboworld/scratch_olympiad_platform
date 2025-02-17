package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UserCore struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	Email          string         `gorm:"not null;"`
	Password       string         `gorm:"not null;"`
	Role           Role           `gorm:"not null;"`
	FullName       string         `gorm:"not null;"`
	FullNameNative string         `gorm:"not null;"`
	CountryID      uint           `gorm:"not null"`
	Country        CountryCore    `gorm:"foreignKey:CountryID"`
	RegionID       *uint          `gorm:"default:null"`
	Region         *RegionCore    `gorm:"foreignKey:RegionID"`
	City           string         `gorm:"not null;"`
	Birthdate      time.Time      `gorm:"not null;"`
	IsActive       bool           `gorm:"not null;default:false;type:boolean;column:is_active"`
}

func (u *UserHTTP) FromCore(userCore UserCore) {
	countryHttp := CountryHTTP{}
	countryHttp.FromCore(userCore.Country)

	var regionHttp *RegionHTTP
	if userCore.Region != nil {
		tmpRegion := RegionHTTP{}
		tmpRegion.FromCore(*userCore.Region)
		regionHttp = &tmpRegion
	}

	u.ID = strconv.Itoa(int(userCore.ID))
	u.CreatedAt = userCore.CreatedAt.Format(time.DateTime)
	u.UpdatedAt = userCore.UpdatedAt.Format(time.DateTime)
	u.Email = userCore.Email
	u.FullName = userCore.FullName
	u.FullNameNative = userCore.FullNameNative
	u.Country = &countryHttp
	u.Region = regionHttp
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
