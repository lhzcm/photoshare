package models

import "time"

type Comment struct {
	Id        int32     `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Pid       int32     `gorm:"not null;index"`                   //动态id
	Cid       int32     `gorm:"not null;default 0"`               //被评论id，0：评论的是动态，大于0：评论的是评论
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
