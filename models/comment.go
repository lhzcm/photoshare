package models

import "time"

type Comment struct {
	Id        int32     `gorm:"primary_key;identity;not null"`
	Pid       int32     `gorm:"not null;index"`             //帖子id
	Cid       int32     `gorm:"not null;default 0"`         //评论id
	Userid    int32     `gorm:"not null"`                   //用户id
	Content   string    `gorm:"not null;type:varchar(512)"` //评论内容
	Praise    int32     `gorm:"not null;default:0"`         //点赞数
	Comments  int32     `gorm:"not null;default:0"`         //评论数
	Writetime time.Time `gorm:"not null;default:getdate();type:datetime"`
}

func (Comment) TableName() string {
	return "t_comment"
}
