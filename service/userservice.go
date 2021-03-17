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
	err = db.GormDB.Where("id = ?", id).First(&user).Error
	redis.Redissetuser(user)
	user.Password = ""
	return
}

//通过用户手机号获取用户信息
func GetUserInfoByPhone(user *User) error {
	if err := db.GormDB.Where("phone = ?", user.Phone).First(&user).Error; err != nil {
		return err
	}
	user.Token = utility.GetTokenById(user.Id, time.Now().Add(time.Hour*24*30))
	if _, err := redis.Redisgetuser(user.Id); err != nil {
		redis.Redissetuser(*user)
	}
	return db.GormDB.Model(&User{Id: user.Id}).Update("token", user.Token).Error
}

//用户登录
func UserLogin(id int32, phone string, password string) (user User, token string, err error) {
	if id <= 0 {
		user.Phone = phone
		if err = GetUserInfoByPhone(&user); err != nil {
			return user, "", errors.New("登录失败，账号有误")
		}
	} else {
		if user, err = GetUserInfoById(id); err != nil {
			return user, "", errors.New("登录失败，账号有误")
		}
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
	var count int64
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

//创建登录/注册验证码PhoneCode
func CreatePhoneCode(phone string, types int) (PhoneCode, error) {
	phonecode := PhoneCode{Phone: phone, Type: types}
	//如果是登录判断该手机号是否注册
	var count int64
	db.GormDB.Model(&User{}).Where("phone = ?", phone).Count(&count)
	if types == 0 && count <= 0 {
		return phonecode, errors.New("该手机号还没有注册，请先去注册账号")
	} else if types == 1 && count > 0 {
		return phonecode, errors.New("该手机号已经注册，请登录！")
	}
	if err := db.GormDB.Create(&phonecode).Error; err != nil {
		return phonecode, errors.New("生成验证码错误")
	}
	return phonecode, nil
}

//验证登录验证码
func ValidatePhoneCode(phone string, code int, types int) (PhoneCode, error) {
	var photoCode PhoneCode
	return photoCode,
		db.GormDB.Where("id = ? and phone = ? and type = ? and expire >= getdate()", code, phone, types).First(&photoCode).Error
}

//获取注册码详情
func GetPhoneCode(id int32) (p PhoneCode, err error) {
	err = db.GormDB.Where("id = ?", id).First(&p).Error
	return
}

//更新PoneCode
func UpdatePhoneCode(id int32, phone string) (int64, error) {
	dbexec := db.GormDB.Model(&PhoneCode{}).Where("id = ? and phone = ? and status = 0 and expire >= getdate()", id, phone).
		Updates(PhoneCode{Status: 1, Updatetime: time.Now()})
	return dbexec.RowsAffected, dbexec.Error
}
