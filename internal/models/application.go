package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ApplicationCore struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	AuthorID   uint           `json:"author_id"`
	User       UserCore       `gorm:"foreignKey:AuthorID"`
	Nomination string         `gorm:"not null" json:"nomination"`
	Link       string         `gorm:"size:255" json:"link"`
	Note       string         `gorm:"size:1024" json:"note"`
}

func (a *ApplicationHTTP) FromCore(application ApplicationCore) {
	a.ID = strconv.Itoa(int(application.ID))
	a.CreatedAt = application.CreatedAt.Format(time.DateTime)
	a.UpdatedAt = application.UpdatedAt.Format(time.DateTime)
	a.AuthorID = strconv.Itoa(int(application.AuthorID))
	a.Nomination = application.Nomination
	a.Link = application.Link
	a.Note = application.Note
}
