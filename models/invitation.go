package models

import "time"

type Invitation struct {
	Id         int32     `gorm:"primary_key;AUTO_INCREMENT"`
	Userid     int32     `gorm:"not null;index"`                           //用户id
	Invitedid  int32     `gorm:"not null;index"`                           //被邀请用户id
	Message    string    `gorm:"null;type:varchar(128)"`                   //邀请信息
	Status     int16     `gorm:"not null;type:smallint;default:0"`         //邀请状态，0：邀请中，1：已邀请，-1：已拒绝，-2：已删除
	Updatetime time.Time `gorm:"not null;type:datetime;default:getdate()"` //更新时间
	Deltime    time.Time `gorm:"not null;type:datetime;default:getdate()"` //用户删除好友的时间
	Deluserid  int32     `gorm:"not null;default:0"`                       //删除操作用户id
	Writetime  time.Time `gorm:"not null;type:datetime;default:getdate()"` //写入时间

	Headimg string `gorm:"->"` //邀请用户头像
	Name    string `gorm:"->"` //邀请用户名称
}

func (Invitation) TableName() string {
	return "t_invitation"
}
