package redis

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	"lkl.photoshare/models"

	"github.com/go-redis/redis/v8"
)

const (
	address  = "127.0.0.1:6379"
	password = ""
)

const (
	userdb   = 0
	cookiedb = 1
	smsdb    = 2
)

var userRDB, smsRDB *redis.Client
var ctx = context.Background()

func init() {
	userRDB = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
	smsRDB = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       2,
	})

	if err := userRDB.Ping(ctx).Err(); err != nil {
		log.Fatalln(err.Error())
	}
	if err := smsRDB.Ping(ctx).Err(); err != nil {
		log.Fatalln(err.Error())
	}
}

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
	user.Cookie = val["Cookie"]
	return
}

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
		"Cookie", user.Cookie).Err()
	if err != nil {
		return
	}
	hours := rand.Intn(5000) + 5000
	err = userRDB.Expire(ctx, strconv.Itoa(int(user.Id)), time.Duration(hours)*time.Hour).Err()
	return
}
