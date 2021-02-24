package service

import (
	"errors"
	db "photoshare/database"
	. "photoshare/models"
	"time"
)

//获取好友列表
func FriendsList(userid int32) (friends []User, err error) {
	query := db.GormDB.Raw("select u.id, u.name, u.headimg "+
		"from "+Friend{}.TableName()+" as f join "+User{}.TableName()+" as u "+
		"on f.friendid = u.id where f.userid = ?", userid)

	if query.Error != nil {
		return nil, query.Error
	}
	err = query.Scan(&friends).Error
	return
	// err = db.GormDB.Where("userid = ?", userid).Find(&friends).Error
	// if err != nil {
	// 	return
	// }
	// for i := 0; i < len(friends); i++ {
	// 	if friends[i].FriendInfo, err = GetUserInfoById(friends[i].Friendid); err != nil {
	// 		return
	// 	}
	// }
	//return
}

//删除好友
func DelFriend(userid int32, friendid int32) error {
	var friend Friend
	if db.GormDB.Where("userid = ? and friendid = ?", userid, friendid).First(&friend).RowsAffected <= 0 {
		return errors.New("删除好友失败，你没有该好友")
	}
	tran := db.GormDB.Begin()
	if tran.Where("userid = ? and friendid = ?", userid, friendid).Delete(&Friend{}).RowsAffected <= 0 {
		tran.Rollback()
		return errors.New("删除好友失败！")
	}
	if tran.Where("userid = ? and friendid = ?", friendid, userid).Delete(&Friend{}).RowsAffected <= 0 {
		tran.Rollback()
		return errors.New("删除好友失败！")
	}
	if tran.Where("id = ?", friend.Iid).Updates(&Invitation{Deltime: time.Now(), Status: -2, Deluserid: userid}).RowsAffected <= 0 {
		tran.Rollback()
		return errors.New("删除好友失败！")
	}
	tran.Commit()
	return nil
}
