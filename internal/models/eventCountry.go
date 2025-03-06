package models

/*
	Таблица связи доступности мероприятия для страны
*/

type EventCountryCore struct {
	ID uint `gorm:"primaryKey"`

	EventID uint
	Event   EventCore `gorm:"foreignKey:EventID;constraint:OnDelete:Cascade;"`

	CountryID uint
	Country   CountryCore `gorm:"foreignKey:CountryID"`
}
