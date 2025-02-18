package models

/*
	Таблица связи доступа пользователя по роли к мероприятию
*/

type EventUserCore struct {
	ID uint `gorm:"primaryKey"`

	EventID uint
	Event   EventCore `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;"`

	UserID    uint
	User      UserCore  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	EventRole EventRole `gorm:"not null"`
}
