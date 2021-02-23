package service

import (
	db "photoshare/database"
	. "photoshare/models"
)

//获取好友列表
func FriendsList(userid int32) (friends []Friend, err error) {
	err = db.GormDB.Where("userid = ?", userid).Find(&friends).Error
	if err != nil {
		return
	}
	for i := 0; i < len(friends); i++ {
		if friends[i].FriendInfo, err = GetUserInfoById(friends[i].Friendid); err != nil {
			return
		}
	}
	return
}
