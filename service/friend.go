package service

import (
	"errors"
	db "photoshare/database"
	. "photoshare/models"
	"time"
)

//用户邀请好友
func Invite(invitation *Invitation) error {
	var count int64 = 0
	db.GormDB.Model(&Invitation{}).Where("userid = ?", invitation.Userid).
		Where("invitedid = ? and status = 0", invitation.Invitedid).Count(&count)
	if count > 0 {
		return errors.New("发送邀请失败！你已经发送邀请了，请等待用户接受")
	}
	db.GormDB.Model(&Friend{}).Where("userid = ?", invitation.Userid).
		Where("friendid = ?", invitation.Invitedid).Count(&count)
	if count > 0 {
		return errors.New("发送邀请失败！你们已经是好友了无需再邀请")
	}
	if db.GormDB.Create(&invitation).RowsAffected <= 0 {
		return errors.New("发送邀请失败！")
	}
	return nil
}

//用户接受或拒绝好友邀请
func InviteAccept(id int32, userid int32, status int16) error {
	var invitation Invitation
	db.GormDB.Where("id = ?", id).Where("invitedid = ? and status = 0", userid).First(&invitation)
	if invitation.Id != id {
		return errors.New("没有找到邀请")
	}
	if status == -1 {
		if db.GormDB.Where("id = ?", id).Updates(&Invitation{Status: status, Updatetime: time.Now()}).RowsAffected <= 0 {
			return errors.New("拒绝失败！请稍后再试")
		}
		return nil
	}

	tran := db.GormDB.Begin()
	if tran.Where("id = ?", id).Updates(&Invitation{Status: status, Updatetime: time.Now()}).RowsAffected <= 0 {
		tran.Rollback()
		return errors.New("接受失败！请稍后再试")
	}
	if tran.Create(&Friend{Userid: invitation.Userid, Friendid: invitation.Invitedid, Iid: id}).RowsAffected <= 0 {
		tran.Rollback()
		return errors.New("添加好友失败")
	}
	if tran.Create(&Friend{Userid: invitation.Invitedid, Friendid: invitation.Userid, Iid: id}).RowsAffected <= 0 {
		tran.Rollback()
		return errors.New("添加好友失败")
	}
	tran.Commit()
	return nil
}

//用户好友邀请列表
func InviteList(userid int32, status int16, page int, pagesize int) (invitation []Invitation, count int64, err error) {
	err = db.GormDB.Model(&Invitation{}).Where("invitedid = ?", userid).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return
	}
	err = db.GormDB.Where("invitedid = ?", userid).Where("status = ?", status).
		Order("id desc").Offset((page - 1) * pagesize).Limit(pagesize).Find(&invitation).Error
	return
}

//获取好友列表
func FriendsList(userid int32) (friends []Friend, err error) {
	err = db.GormDB.Where("userid = ?", userid).Find(&friends).Error
	return
}
