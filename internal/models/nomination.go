package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type NominationCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"unique;size:255"`
	MinAge    uint
	MaxAge    uint
}

func (n *NominationHTTP) FromCore(nomination NominationCore) {
	n.ID = strconv.Itoa(int(nomination.ID))
	n.CreatedAt = nomination.CreatedAt.Format(time.DateTime)
	n.UpdatedAt = nomination.UpdatedAt.Format(time.DateTime)
	n.Name = nomination.Name
	n.MinAge = int(nomination.MinAge)
	n.MaxAge = int(nomination.MaxAge)
}

func FromNominationsCore(nominationsCore []NominationCore) (nominationsHttp []*NominationHTTP) {
	for _, nominationCore := range nominationsCore {
		var tmpNominationHttp NominationHTTP
		tmpNominationHttp.FromCore(nominationCore)
		nominationsHttp = append(nominationsHttp, &tmpNominationHttp)
	}
	return
}
