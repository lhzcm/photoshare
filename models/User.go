package models

import (
	"time"
)

type User struct {
	Id          int32     `gorm:"primary_key;AUTO_INCREMENT"`
	Name        string    `gorm:"not null"`
	Headimg     string    `gorm:"not null"`
	Phone       string    `gorm:"not null;type:char(11)"`
	City        int32     `gorm:"not null"`
	Brithday    time.Time `gorm:"not null;type:date"`
	Ismale      bool      `gorm:"not null"`
	Password    string    `gorm:"not null"`
	Updatetime  time.Time `gorm:"not null;type:datetime;default:getdate()"`
	Writetime   time.Time `gorm:"not null;type:datetime;default:getdate()"`
	Token       string    `gorm:"not null;type:varchar(128)"`
	Notreadnums int       `gorm:"->"` //好友未读消息
	Code        int       `gorm:"-"`  //验证码
}

//定义表名
func (User) TableName() string {
	return "t_users"
}
