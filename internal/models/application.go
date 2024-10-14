package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ApplicationCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	AuthorID  uint
	Author    UserCore `gorm:"foreignKey:AuthorID"`

	Nomination                    string `gorm:"not null"`
	AlgorithmicTaskLink           string `gorm:"size:255"`
	AlgorithmicTaskFile           string `gorm:"size:255"`
	CreativeTaskLink              string `gorm:"size:255"`
	CreativeTaskFile              string `gorm:"size:255"`
	EngineeringTaskFile           string `gorm:"size:255"`
	EngineeringTaskCloudLink      string `gorm:"size:255"`
	EngineeringTaskVideo          string `gorm:"size:255"`
	EngineeringTaskVideoCloudLink string `gorm:"size:255"`
	Note                          string `gorm:"size:1024"`
}

func (a *ApplicationHTTP) FromCore(application ApplicationCore) {
	a.ID = strconv.Itoa(int(application.ID))
	a.CreatedAt = application.CreatedAt.Format(time.DateTime)
	a.UpdatedAt = application.UpdatedAt.Format(time.DateTime)
	a.AuthorID = strconv.Itoa(int(application.AuthorID))
	a.Nomination = application.Nomination
	a.AlgorithmicTaskLink = application.AlgorithmicTaskLink
	a.AlgorithmicTaskFile = application.AlgorithmicTaskFile
	a.CreativeTaskLink = application.CreativeTaskLink
	a.CreativeTaskFile = application.CreativeTaskFile
	a.EngineeringTaskFile = application.EngineeringTaskFile
	a.EngineeringTaskCloudLink = application.EngineeringTaskCloudLink
	a.EngineeringTaskVideo = application.EngineeringTaskVideo
	a.EngineeringTaskVideoCloudLink = application.EngineeringTaskVideoCloudLink
	a.Note = application.Note
}

func FromApplicationsCore(applicationsCore []ApplicationCore) (applicationsHttp []*ApplicationHTTP) {
	for _, applicationCore := range applicationsCore {
		var tmpApplicationHttp ApplicationHTTP
		tmpApplicationHttp.FromCore(applicationCore)
		applicationsHttp = append(applicationsHttp, &tmpApplicationHttp)
	}
	return
}
