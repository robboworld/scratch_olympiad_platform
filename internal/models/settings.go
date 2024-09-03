package models

type SettingsCore struct {
	ID               uint `gorm:"primaryKey" json:"id"`
	ActivationByLink bool `gorm:"not null;default:true;type:boolean;column:activation_by_link"`
}

func (p *Settings) FromCore(settingsCore SettingsCore) {
	p.ActivationByLink = settingsCore.ActivationByLink
}
