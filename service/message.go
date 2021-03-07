package service

import (
	"errors"
	db "photoshare/database"
	"photoshare/models"
	. "photoshare/models"
)

func SendMessage(msg *Message) error {
	var count int64
	db.GormDB.Model(&models.Friend{}).Where("userid = ? and friendid = ?", msg.Senderid, msg.Receiverid).Count(&count)
	if count == 0 {
		return errors.New("消息发送失败！该用户不是你的好友")
	}
	if err := db.GormDB.Create(msg).Error; err != nil {
		return errors.New("消息发送失败！")
	}
	return nil
}
