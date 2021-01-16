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
	Type      uint8     `gorm:"not null;type:tinyint;default:0"` //类型，0：只有自己看，1：好友可看，3：所有人可以看
	Lng       float64   `gorm:"null"`
	Lat       float64   `gorm:"null"`
	Position  string    `gorm:"null;type:varchar(64)"`
	Status    int16     `gorm:"not null;type:smallint;default:0"` //状态，0：正常，-1：已删除
	Dateflag  time.Time `gorm:"not null;type:date;default:getdate()"`
	Writetime time.Time `gorm:"not null;default:getdate();type:datetime"`
	Imgs      string    `gorm:"-"`
	Photos    []Photo   `gorm:"-"`
}

func (Publish) TableName() string {
	return "t_publish"
}
