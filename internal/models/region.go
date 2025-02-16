package models

import "strconv"

type RegionCore struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;size:255"`

	CountryID uint
	Country   CountryCore `gorm:"foreignKey:CountryID;constraint:OnDelete:CASCADE;"`
}

func (r *RegionHTTP) FromCore(region RegionCore) {
	r.ID = strconv.Itoa(int(region.ID))
	r.Name = region.Name
	r.CountryID = strconv.Itoa(int(region.CountryID))
}

func FromRegionsCore(regionsCore []RegionCore) (regionsHttp []*RegionHTTP) {
	for _, regionCore := range regionsCore {
		var tmpRegionHttp RegionHTTP
		tmpRegionHttp.FromCore(regionCore)
		regionsHttp = append(regionsHttp, &tmpRegionHttp)
	}
	return
}
