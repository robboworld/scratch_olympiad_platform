package models

import (
	"strconv"
)

type CountryCore struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;size:255"`
}

func (c *CountryHTTP) FromCore(country CountryCore) {
	c.ID = strconv.Itoa(int(country.ID))
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
