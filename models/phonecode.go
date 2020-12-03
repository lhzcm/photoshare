package models

import (
	"time"
)

type PhoneCode struct {
	Id         int32     `gorm:"primary_key;AUTO_INCREMENT" json:"Id,string"`
	Phone      string    `gorm:"null;type:char(11);default:null"`
	Writetime  time.Time `gorm:"not null;default:getdate();type:datetime"`
	Updatetime time.Time `gorm:"null;type:datetime;default:null"`
}

func (PhoneCode) TableName() string {
	return "t_phone_code"
}
