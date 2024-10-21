package entities

import "time"

type States struct {
	Id        uint       `gorm:"primary_key:auto_increment" json:"id"`
	Name      string     `gorm:"type:varchar(100);not null" json:"name"`
	ZipCode   string     `gorm:"type:varchar(100);null" json:"zip_code"`
	CityId    uint       `gorm:"null" json:"cityid"`
	City      City       `gorm:"foreignkey:CityId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"city"`
	Active    bool       `gorm:"type:boolean" json:"active"`
	CreatedAt time.Time  `gorm:"<-:created_at" json:"created"`
	UpdatedAt *time.Time `gorm:"type:TIMESTAMP(6)"  json:"updated"`
}
