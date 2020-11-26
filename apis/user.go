package apis

import (
	"log"
	"net/http"
	"strconv"

	"photoshare/config"
	"photoshare/middleware"
	. "photoshare/models"
	"photoshare/service"
	"photoshare/utility"

	"github.com/gin-gonic/gin"
)

//route配置
func (router *Router) UserRouteRegister() {
	router.GET("/test", Test)
	router.POST("/user", UserRegister)
	router.POST("/user/login", UserLogin)
	router.Use(middleware.UserValidate).GET("/user/info", UserInfo)
}

//获取登录用户信息
func UserInfo(c *gin.Context) {
	user, isexists := c.Get("user")
	if !isexists {
		c.JSON(http.StatusOK, Fail("没有用户信息"))
	}
	c.JSON(http.StatusOK, Success(user, "成功获取用户信息"))
}

//通过用户id获取用户信息
func UserInfoById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Fail("请求的参数有误"+err.Error()))
		return
	}
	var user User
	if user, err = service.GetUserInfoById(int32(id)); err != nil {
		c.JSON(http.StatusOK, Fail("没有找到用户信息"))
		return
	}
	c.JSON(http.StatusOK, Success(user, "成功获取用户信息"))
}

//用户注册添加
func UserRegister(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusOK, Fail("提交的参数有误"))
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
	if len(user.Password) < 6 {
		c.JSON(http.StatusOK, Fail("用户密码必须大于6位数"))
		return
	}

	exists, err := user.IsPhoneExits()
	if err != nil {
		c.JSON(http.StatusOK, Fail("数据请求出错"))
		log.Println(err)
		return
	}
	if exists {
		c.JSON(http.StatusOK, Fail("当前手机号码已经注册了"))
		return
	}
	user.Password = utility.EncryptPassword(user.Password)
	_, err = user.Insert()
	if err != nil {
		c.JSON(http.StatusOK, Fail("注册失败，请稍后再试"))
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, Success(user, "注册成功"))
}

//用户登录
func UserLogin(c *gin.Context) {
	var user User
	var err error
	var token string
	if err = c.BindJSON(&user); err != nil {
		c.JSON(http.StatusOK, Fail("提交的参数有误"))
		log.Println(err)
		return
	}
	if user, token, err = service.UserLogin(user.Id, user.Password); err != nil {
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	c.SetCookie("token", token, 60*24*30, "/", c.Request.URL.Host, false, false)
	c.JSON(http.StatusOK, Success(user, "登录成功"))
}

//test
func Test(c *gin.Context) {
	c.JSON(http.StatusOK, Success(config.Configs, "成功"))
}
