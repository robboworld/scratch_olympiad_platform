package models

/*
	Таблица связи доступности мероприятия для региона
*/

type EventRegionCore struct {
	ID uint `gorm:"primaryKey"`

	EventID uint
	Event   EventCore `gorm:"foreignKey:EventID;constraint:OnDelete:Cascade;"`

	RegionID uint
	Region   RegionCore `gorm:"foreignKey:RegionID"`
}
