package models

import "time"

type Friend struct {
	Userid    int32     `gorm:"primary_key"`                              //用户id
	Friendid  int32     `gorm:"primary_key"`                              //好友id
	Iid       int32     `gorm:"not null"`                                 //邀请id
	Writetime time.Time `gorm:"not null;type:datetime;default:getdate()"` //添加时间
}

func (Friend) TableName() string {
	return "t_friend"
}
