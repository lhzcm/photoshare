package service

import (
	"errors"
	db "photoshare/database"
	"photoshare/models"
	. "photoshare/models"
)

//发送消息
func SendMessage(msg *Message) error {
	var friend Friend
	db.GormDB.Model(&models.Friend{}).Where("userid = ? and friendid = ?", msg.Receiverid, msg.Senderid).First(&friend)
	if friend.Userid <= 0 {
		return errors.New("消息发送失败！该用户不是你的好友")
	}
	if err := db.GormDB.Create(msg).Error; err != nil {
		return errors.New("消息发送失败！")
	}
	db.GormDB.Model(&models.Friend{}).Where("userid = ? and friendid = ?", msg.Receiverid, msg.Senderid).Update("notreadnums", friend.Notreadnums+1)
	return nil
}

//获取消息列表
func GetMessageList(userid int, senderid int, curid int) (msgs []models.Message, err error) {
	err = db.GormDB.Where("((senderid = ? and receiverid = ?) or (senderid = ? and receiverid = ?)) and id < ?",
		userid, senderid, senderid, userid, curid).Order("id desc").Limit(20).Find(&msgs).Error
	//删除未读消息数
	if curid == 1000000000 {
		UpdateMessageNotreadnums(userid, senderid)
	}
	return
}

//更新未读消息数
func UpdateMessageNotreadnums(userid int, friendid int) error {
	return db.GormDB.Model(&Friend{}).Where("userid = ? and friendid = ? and notreadnums > 0",
		userid, friendid).Update("notreadnums", 0).Error
}
