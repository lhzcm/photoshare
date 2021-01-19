package models

import "time"

type Praise struct {
	Rid       int32     `gorm:"not null;primary_key"`                         //关联被点赞id
	Userid    int32     `gorm:"not null;primary_key"`                         //用户id
	Type      int8      `gorm:"not null;primary_key;type:tinyint;default:0;"` //被点赞类型，0：动态，1：评论
	Writetime time.Time `gorm:"not null;type:datetime;default:getdate()"`
}

func (Praise) TableName() string {
	return "t_praise"
}
