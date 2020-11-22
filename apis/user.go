package apis

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	. "lkl.photoshare/models"
	"lkl.photoshare/redis"
)

//route配置
func (router *Router) UserRouteRegister() {
	router.GET("/user/:id", UserInfo)
	router.POST("/user", UserAdd)
}

//通过用户id获取用户信息
func UserInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Fail("请求的参数有误"))
		log.Println(err)
		return
	}
	u := User{Id: int32(id)}
	var user User
	if user, err = redis.Redisgetuser(u.Id); err == nil {
		c.JSON(http.StatusOK, Success(user, "数据请求成功"))
		return
	}
	if err = u.GetFirst(); err != nil {
		c.JSON(http.StatusOK, Fail("没有找到数据"))
		log.Println(err)
		return
	}
	redis.Redissetuser(u)
	c.JSON(http.StatusOK, Success(u, "数据请求成功"))
}

//用户注册添加
func UserAdd(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusOK, Fail("提交的参数有误"))
		log.Println(err)
		return
	}
	exists, err := user.IsPhoneExits()
	if err != nil {
		c.JSON(http.StatusOK, Fail("数据请求出错"))
		log.Println(err)
		return
	}
	if len(user.Name) < 3 && len(user.Name) > 16 {
		c.JSON(http.StatusOK, Fail("用户名有误，请输入3-16个字符的用户名"))
		return
	}
	if len(user.Phone) != 11 {
		c.JSON(http.StatusOK, Fail("手机号码有误"))
		return
	}
	if exists {
		c.JSON(http.StatusOK, Fail("当前手机号码已经注册了"))
		return
	}

	_, err = user.Insert()
	if err != nil {
		c.JSON(http.StatusOK, Fail("注册失败，请稍后再试"))
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, Success(user, "注册成功"))
}
