package redis

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	"photoshare/config"
	"photoshare/models"

	"github.com/go-redis/redis/v8"
)

var userRDB, msgRDB *redis.Client
var ctx = context.Background()

func init() {
	//用户信息redis库
	dbconfig := config.Configs.Redis.User
	userRDB = redis.NewClient(&redis.Options{
		Addr:     dbconfig.Address,
		Password: dbconfig.Password,
		DB:       dbconfig.Dbnum,
	})
	//用户聊天消息redis库
	dbconfig = config.Configs.Redis.Msg
	msgRDB = redis.NewClient(&redis.Options{
		Addr:     dbconfig.Address,
		Password: dbconfig.Password,
		DB:       dbconfig.Dbnum,
	})

	if err := userRDB.Ping(ctx).Err(); err != nil {
		log.Fatalln(err.Error())
	}
	if err := msgRDB.Ping(ctx).Err(); err != nil {
		log.Fatalln(err.Error())
	}
}

//获取用户缓存
func Redisgetuser(Id int32) (user models.User, err error) {
	var val map[string]string
	if val, err = userRDB.HGetAll(ctx, strconv.Itoa(int(Id))).Result(); err != nil {
		return
	}
	if len(val) == 0 {
		err = errors.New("未找到用户信息")
		return
	}
	var inttemp int
	inttemp, _ = strconv.Atoi(val["Id"])
	user.Id = int32(inttemp)
	user.Name = val["Name"]
	user.Headimg = val["Headimg"]
	user.Phone = val["Phone"]
	inttemp, _ = strconv.Atoi(val["City"])
	user.City = int32(inttemp)
	timestamp, _ := strconv.Atoi(val["Brithday"])
	user.Brithday = time.Unix(int64(timestamp), 0)
	if val["Ismale"] == "true" {
		user.Ismale = true
	}
	user.Password = val["Password"]
	timestamp, _ = strconv.Atoi(val["Updatetime"])
	user.Updatetime = time.Unix(int64(timestamp), 0)
	timestamp, _ = strconv.Atoi(val["Writetime"])
	user.Writetime = time.Unix(int64(timestamp), 0)
	user.Token = val["Token"]
	return
}

//用户信息缓存
func Redissetuser(user models.User) (err error) {
	err = userRDB.HMSet(ctx, strconv.Itoa(int(user.Id)),
		"Id", user.Id,
		"Name", user.Name,
		"Headimg", user.Headimg,
		"Phone", user.Phone,
		"City", user.City,
		"Brithday", user.Brithday.Unix(),
		"Ismale", user.Ismale,
		"Password", user.Password,
		"Updatetime", user.Updatetime.Unix(),
		"Writetime", user.Writetime.Unix(),
		"Token", user.Token).Err()
	if err != nil {
		return
	}
	hours := rand.Intn(5000) + 5000
	err = userRDB.Expire(ctx, strconv.Itoa(int(user.Id)), time.Duration(hours)*time.Hour).Err()
	return
}

//token缓存
func RedissetToken(id int32, token string) (err error) {
	err = userRDB.HMSet(ctx, strconv.FormatInt(int64(id), 10),
		"Token", token).Err()
	if err != nil {
		return
	}
	return
}

//聊天消息缓存
func RedisAddMsg(receiverid int32, msg []byte) (err error) {
	key := strconv.Itoa(int(receiverid))
	if err = msgRDB.LPush(ctx, key, msg).Err(); err != nil {
		return
	}
	return msgRDB.Expire(ctx, key, time.Hour*24).Err()
}

//获取聊天信息缓存
func RedisGetMsg(receiverid int32) ([]byte, error) {
	str, err := msgRDB.RPop(ctx, strconv.Itoa(int(receiverid))).Bytes()
	return str, err
}

//获取缓存聊天数
func RedisGetMsgCount(receiverid int32) int64 {
	return msgRDB.LLen(ctx, strconv.Itoa(int(receiverid))).Val()
}
