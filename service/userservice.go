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
	user.Id = id
	if err = user.GetFirst(); err != nil {
		log.Println(err)
		return
	}
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
