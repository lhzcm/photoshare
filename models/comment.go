package models

import "time"

type Comment struct {
	Id        int32     `gorm:"primary_key;identity;not null"`
	Rid       int32     `gorm:"not null;index"`                   //评论关联id
	Type      int8      `gorm:"not null;type:tinyint;default 0"`  //评论类型，0：动态评论， 1：用户评论的评论
	Userid    int32     `gorm:"not null"`                         //用户id
	Content   string    `gorm:"not null;type:varchar(512)"`       //评论内容
	Praise    int32     `gorm:"not null;default:0"`               //点赞数
	Comments  int32     `gorm:"not null;default:0"`               //评论数
	Status    int16     `gorm:"not null;type:smallint;default:0"` //状态， 0：正常，-1已删除
	Writetime time.Time `gorm:"not null;default:getdate();type:datetime"`
}

func (Comment) TableName() string {
	return "t_comment"
}
