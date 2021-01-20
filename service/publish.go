package service

import (
	"errors"
	db "photoshare/database"
	. "photoshare/models"
	"strconv"
)

//上传图片
func AddUploadImg(p *Photo) error {
	return db.GormDB.Create(p).Error
}

//判断图片是否是用户上传的
func PhotoIsUser(ids []int, userid int) int {
	if len(ids) == 0 {
		return 0
	}

	strs := make([]string, len(ids))
	for i, item := range ids {
		strs[i] = strconv.Itoa(item)
	}
	var count int = 0
	db.GormDB.Model(&Photo{}).Where("id in (?)", ids).
		Where("userid = ?", userid).Count(&count)
	return count
}

//保存用户发表的动态
func SavePublish(p *Publish, ids []int) error {
	tx := db.GormDB.Begin()
	if err := tx.Create(p).Error; err != nil {
		tx.Rollback()
		return err
	}
	if p.Id <= 0 {
		return errors.New("发表失败")
	}
	if len(ids) > 0 && tx.Model(&Photo{}). //tx.Table(Photo{}.TableName()).
						Where("id in (?)", ids).Where("pid = 0").Updates(&Photo{Pid: p.Id}).RowsAffected != int64(len(ids)) {
		tx.Rollback()
		return errors.New("发表失败，更新图片失败")
	}
	tx.Commit()
	return nil
}

//删除动态
func DeletePublish(id int32, userid int32) error {
	exec := db.GormDB.Model(&Publish{}).Where("id = ?", id).
		Where("userid = ?", userid).Update(&Publish{Status: -1})
	if exec.Error != nil {
		return exec.Error
	}
	if exec.RowsAffected <= 0 {
		return errors.New("没有找到可以删除的动态")
	}
	return nil
}

//获取用户动态列表
func GetPublishList(userid int, page int, pagesize int) (publishs []Publish, total int, err error) {
	var photos []Photo

	db.GormDB.Model(&Publish{}).Where("userid = ? and status = 0", userid).Count(&total)

	if db.GormDB.Where("userid = ? and status = 0", userid).Offset((page-1)*pagesize).
		Limit(pagesize).Order("id desc").Find(&publishs).Error != nil {
		return publishs, total, errors.New("获取动态失败")
	}

	for i := 0; i < len(publishs); i++ {
		if db.GormDB.Where("pid = ?", publishs[i].Id).Order("id").Find(&photos).Error != nil {
			return publishs, total, errors.New("获取动态图片失败")
		}
		publishs[i].Photos = photos
	}

	return publishs, total, nil
}

//用户点赞
func PublishPraise(userid int32, id int, ptype int) error {
	var publish Publish
	var comment Comment
	if ptype == 0 { //动态点赞
		db.GormDB.Where("id = ?", id).Where("status = 0").First(&publish)
		if publish.Id == 0 {
			return errors.New("点赞失败，没有找到对应的动态")
		}
		tran := db.GormDB.Begin()
		if tran.Model(&Publish{}).Where("id = ?", id).Update(&Publish{Praise: publish.Praise + 1}).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("点赞失败")
		}
		if tran.Create(&Praise{Rid: publish.Id, Userid: userid, Type: 0}).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("点赞失败")
		}
		tran.Commit()
	} else { //评论点赞
		db.GormDB.Where("id = ?", id).Where("status = 0").First(&comment)
		if comment.Id == 0 {
			return errors.New("点赞失败，没有找到对应的评论")
		}
		tran := db.GormDB.Begin()
		if tran.Model(&Comment{}).Where("id = ?", id).Update(&Comment{Praise: comment.Praise + 1}).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("点赞失败")
		}
		if tran.Create(&Praise{Rid: comment.Id, Userid: userid, Type: 1}).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("点赞失败")
		}
		tran.Commit()
	}

	return nil
}

//用户取消点赞
func PublishUnPraise(userid int32, id int, ptype int) error {
	var publish Publish
	var comment Comment
	if ptype == 0 { //动态点赞
		db.GormDB.Where("id = ?", id).Where("status = 0").First(&publish)
		if publish.Id == 0 {
			return errors.New("取消点赞失败，没有找到对应的动态")
		}
		tran := db.GormDB.Begin()
		if tran.Model(&Publish{}).Where("id = ?", id).Update("praise", publish.Praise-1).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("取消点赞失败")
		}
		if tran.Where("rid = ?", id).Where("userid = ?", userid).Where("type = ?", 0).Delete(&Praise{}).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("取消点赞失败")
		}
		tran.Commit()
	} else { //评论点赞
		db.GormDB.Where("id = ?", id).Where("status = 0").First(&comment)
		if comment.Id == 0 {
			return errors.New("取消点赞失败，没有找到对应的评论")
		}
		tran := db.GormDB.Begin()
		if tran.Model(&Comment{}).Where("id = ?", id).Update("praise", comment.Praise-1).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("取消点赞失败")
		}

		if tran.Where("rid = ?", id).Where("userid = ?", userid).Where("type = ?", 1).Delete(&Praise{}).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("取消点赞失败")
		}
		tran.Commit()
	}

	return nil
}

func PublishComment(c *Comment) error {
	var publish Publish
	var comment Comment
	if c.Cid == 0 { //评论动态
		db.GormDB.Where("id = ?", c.Pid).Where("status = 0").First(&publish)
		if publish.Id == 0 {
			return errors.New("评论失败，没有找到对应的动态")
		}
		tran := db.GormDB.Begin()
		if tran.Model(&Publish{}).Where("id = ?", c.Pid).Update("comments", publish.Comments+1).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("评论失败")
		}
		if tran.Create(c).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("评论失败")
		}
		tran.Commit()
	} else { //评论评论
		db.GormDB.Where("id = ?", c.Cid).Where("pid = ?", c.Pid).Where("status = 0").First(&comment)
		if comment.Id == 0 {
			return errors.New("评论失败，没有找到对应的评论")
		}
		tran := db.GormDB.Begin()
		if tran.Model(&Comment{}).Where("id = ?", c.Cid).Update("comments", comment.Comments+1).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("评论失败")
		}
		if tran.Create(c).RowsAffected == 0 {
			tran.Rollback()
			return errors.New("评论失败")
		}
		tran.Commit()
	}

	return nil

}
