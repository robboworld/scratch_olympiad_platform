package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type CountryCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"unique;size:255"`
}

func (c *CountryHTTP) FromCore(country CountryCore) {
	c.ID = strconv.Itoa(int(country.ID))
	c.CreatedAt = country.CreatedAt.Format(time.DateTime)
	c.UpdatedAt = country.UpdatedAt.Format(time.DateTime)
	c.Name = country.Name
}

func FromCountriesCore(countriesCore []CountryCore) (countriesHttp []*CountryHTTP) {
	for _, countryCore := range countriesCore {
		var tmpCountryHttp CountryHTTP
		tmpCountryHttp.FromCore(countryCore)
		countriesHttp = append(countriesHttp, &tmpCountryHttp)
	}
	return
}
