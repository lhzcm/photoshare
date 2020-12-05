package models

import (
	"time"
)

type Publish struct {
	Id        int32     `gorm:"primary_key;AUTO_INCREMENT"`
	Userid    int32     `gorm:"not null;index"`
	Content   string    `gorm:"not null;type:varchar(2048)"`
	Praise    int32     `gorm:"not null;default:0"`
	Comments  int32     `gorm:"not null;default:0"`
	Lng       float64   `gorm:"null"`
	Lat       float64   `gorm:"null"`
	Position  string    `gorm:"null;type:varchar(64)"`
	Dateflag  time.Time `gorm:"not null;type:date;default:getdate()"`
	Writetime time.Time `gorm:"not null;default:getdate();type:datetime"`
}

func (Publish) TableName() string {
	return "t_publish"
}
