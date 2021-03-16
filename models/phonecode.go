package models

import (
	"time"
)

type PhoneCode struct {
	Id         int32     `gorm:"primary_key;AUTO_INCREMENT"`                                   //验证码
	Phone      string    `gorm:"not null;type:char(11);default:null"`                          //手机号
	Type       int       `gorm:"not null;type:int;default:0"`                                  //验证码类型，0：登录验证， 1：注册验证
	Status     int       `gorm:"not null;type:int;default:0"`                                  //验证码状态，0：未验证，1：已验证
	Expire     time.Time `gorm:"not null;index;type:datetime;default:dateadd(hh,2,getdate())"` //验证码过期时间
	Writetime  time.Time `gorm:"not null;default:getdate();type:datetime"`                     //添加时间
	Updatetime time.Time `gorm:"null;type:datetime;default:getdate()"`                         //更新时间
}

func (PhoneCode) TableName() string {
	return "t_phone_code"
}
