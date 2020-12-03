package service

import (
	"errors"
	"log"
	db "photoshare/database"
	. "photoshare/models"
	"photoshare/redis"
	"photoshare/utility"
	"time"
)

//通过用户id获取用户信息（优先从redis里面获取）
func GetUserInfoById(id int32) (user User, err error) {
	if user, err = redis.Redisgetuser(id); err == nil {
		return
	}
	db.GormDB.Where("id = ?", id).First(&user)
	redis.Redissetuser(user)
	return
}

//用户登录
func UserLogin(id int32, password string) (user User, token string, err error) {
	if user, err = GetUserInfoById(id); err != nil {
		return user, "", errors.New("登录失败，账号有误")
	}
	if user.Password != utility.EncryptPassword(password) {
		return user, "", errors.New("密码错误")
	}
	token = utility.GetTokenById(user.Id, time.Now().Add(time.Hour*24*30))
	db.GormDB.Model(&User{Id: id}).Update("token", token)
	redis.RedissetToken(id, token)
	return
}

//用注册
func UserRegister(user *User) error {
	var count int
	if err := db.GormDB.Where("phone = ?", user.Phone).Find(&User{}).Count(&count).Error; err != nil {
		log.Println(err)
		return errors.New("系统错误")
	}
	if count > 0 {
		return errors.New("该手机号已被注册")
	}
	if err := db.GormDB.Create(user).Error; err != nil {
		return errors.New("注册失败")
	}
	return nil
}

//创建注册码PhoneCode
func CreatePhoneCode() (int32, error) {
	phonecode := PhoneCode{}
	err := db.GormDB.Create(&phonecode).Error
	return phonecode.Id, err
}

//获取注册码详情
func GetPhoneCode(id int32) (p PhoneCode, err error) {
	err = db.GormDB.Where("id = ?", id).First(&p).Error
	return
}

//更新PoneCode
func UpdatePhoneCode(id int32, phone string) (int64, error) {
	dbexec := db.GormDB.Model(&PhoneCode{}).Where("id = ?", id).Where("phone is null").
		Update(PhoneCode{Phone: phone, Updatetime: time.Now()})
	return dbexec.RowsAffected, dbexec.Error
}
