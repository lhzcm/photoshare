package utility

import (
	"crypto/md5"
	"errors"
	"fmt"
	"photoshare/config"
	"strconv"
	"strings"
	"time"
)

var tokenSecret = config.Configs.Token.Secret

const passwordSecret = "password_photoshare"

//通过id和过期时间加密成Token
func GetTokenById(id int32, expire time.Time) string {
	idStr := strconv.FormatInt(int64(id), 10)
	timeStr := strconv.FormatInt(expire.Unix(), 10)

	bytes := []byte(idStr + tokenSecret + timeStr)
	has := md5.Sum(bytes)
	return idStr + "-" + fmt.Sprintf("%x", has) + "-" + timeStr
}

//通过Token解密成id和过期时间
func GetIdByToken(token string) (int32, error) {
	strs := strings.Split(token, "-")
	if len(strs) != 3 {
		return 0, errors.New("无效的token")
	}

	id, err := strconv.ParseInt(strs[0], 10, 32)
	if err != nil {
		return 0, err
	}
	var timespan int64
	timespan, err = strconv.ParseInt(strs[2], 10, 64)
	if err != nil {
		return 0, err
	}
	if timespan < time.Now().Unix() {
		return 0, errors.New("token已经过期")
	}
	if GetTokenById(int32(id), time.Unix(timespan, 0)) != token {
		return 0, errors.New("token以破坏")
	}
	return int32(id), err
}

//加密密码
func EncryptPassword(pwd string) string {
	bytes := []byte(pwd + passwordSecret)
	md5Bytes := md5.Sum(bytes)
	return fmt.Sprintf("%x", md5Bytes)
}
