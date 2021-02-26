package service

import (
	"errors"
	db "photoshare/database"
	. "photoshare/models"
)

func SendMessage(msg *Message) error {
	if err := db.GormDB.Create(msg).Error; err != nil {
		return errors.New("消息发送失败！")
	}
	return nil
}
