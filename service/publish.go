package service

import (
	db "photoshare/database"
	. "photoshare/models"
)

func AddUploadImg(p *Photo) error {
	return db.GormDB.Create(p).Error
}
