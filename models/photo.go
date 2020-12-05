package models

import (
	"time"
)

type Photo struct {
	Id        int32     `gorm:"primary key;identity"`
	Userid    int32     `gorm:"not null"`
	Pid       int32     `gorm:"not null;index;default:0"`
	Imgurl    string    `gorm:"not null;type:varchar(256)"`
	Info      string    `gorm:"null;type:varchar(1024)"`
	Writetime time.Time `gorm:"not null;default:getdate();type:datetime"`
}
