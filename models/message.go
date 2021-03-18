package models

import "time"

type Message struct {
	Id         int32     `gorm:"primary_key;AUTO_INCREMENT;type:int"`
	Senderid   int32     `gorm:"not null;type:int;index:idx_t_message_sr"` //发送者id
	Receiverid int32     `gorm:"not null;type:int;index:idx_t_message_sr"` //接收者id
	Content    string    `gorm:"not null;type:nvarchar(1024)"`             //消息内容
	Type       int16     `gorm:"not null;default:1;type:tinyint"`          //消息类型 1:文字信息，2：图片， 3：文件
	Status     int16     `gorm:"not null;default:0;type:tinyint"`          //消息状态 0：正常， -1：撤回， -2：删除
	Writetime  time.Time `gorm:"not null;default:getdate();type:datetime"` //发送时间
}

func (Message) TableName() string {
	return "t_message"
}
