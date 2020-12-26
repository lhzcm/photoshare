package apis

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"

	"photoshare/middleware"
	. "photoshare/models"
	"photoshare/service"
	"photoshare/utility"

	"github.com/gin-gonic/gin"
)

//route配置
func (router *Router) UserRouteRegister() {
	router.Use(gin.Logger())
	router.GET("/user/code", GetPhoneCode)
	router.GET("/user/code/:id", PhoneCodeInfo)
	router.GET("/user/callback", CallBackPhoneCode)
	router.POST("/user", UserRegister)
	router.POST("/user/login", UserLogin)
	router.GET("/user/info", middleware.UserValidate, UserInfo)
}

//获取登录用户信息
func UserInfo(c *gin.Context) {
	user, isexists := c.Get("user")
	if !isexists {
		c.JSON(http.StatusOK, Fail("没有用户信息"))
	}
	c.JSON(http.StatusOK, Success(user, "成功获取用户信息"))
}

//获取用户信息
func GetUserInfo(c *gin.Context) User {
	if user, isexists := c.Get("user"); isexists {
		return user.(User)
	}
	return User{}
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

	user.Password = utility.EncryptPassword(user.Password)

	if err := service.UserRegister(&user); err != nil {
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	if _, token, err := service.UserLogin(user.Id, user.Password); err == nil {
		c.SetCookie("token", token, 60*24*30, "/", c.Request.URL.Host, false, false)
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

//用户注册返回code
func GetPhoneCode(c *gin.Context) {
	if code, err := service.CreatePhoneCode(); err == nil {
		c.JSON(http.StatusOK, Success(code, "成功"))
		return
	}
	c.JSON(http.StatusOK, Fail("系统错误"))
}

//第三方回调Phone注册码
func CallBackPhoneCode(c *gin.Context) {
	phone := c.Query("mobile")
	result := KeliResult{Version: "1.0"}

	code, err := strconv.ParseInt(c.Query("content"), 10, 32)
	if err != nil {
		result.ResultInfo = "返回code码有误"
		xmlhead := []byte(xml.Header)
		xmlcontent, _ := xml.MarshalIndent(result, "", "    ")
		xmlbytes := append(xmlhead, xmlcontent...)

		c.String(http.StatusOK, string(xmlbytes))
		return
	}

	if len(phone) != 11 {
		result.ResultInfo = "手机号码有误"
		xmlhead := []byte(xml.Header)
		xmlcontent, _ := xml.MarshalIndent(result, "", "    ")
		xmlbytes := append(xmlhead, xmlcontent...)
		//c.XML(http.StatusOK, result)
		c.String(http.StatusOK, string(xmlbytes))
		return
	}

	var rowcount int64
	if rowcount, err = service.UpdatePhoneCode(int32(code), phone); err != nil {
		result.ResultInfo = err.Error()
		xmlhead := []byte(xml.Header)
		xmlcontent, _ := xml.MarshalIndent(result, "", "    ")
		xmlbytes := append(xmlhead, xmlcontent...)
		//c.XML(http.StatusOK, result)
		c.String(http.StatusOK, string(xmlbytes))
		return
	}
	if rowcount <= 0 {
		result.ResultInfo = "更新失败"
		xmlhead := []byte(xml.Header)
		xmlcontent, _ := xml.MarshalIndent(result, "", "    ")
		xmlbytes := append(xmlhead, xmlcontent...)
		//c.XML(http.StatusOK, result)
		c.String(http.StatusOK, string(xmlbytes))
		return
	}

	result.Result = 1
	result.ResultInfo = "注册码发送成功"
	xmlhead := []byte(xml.Header)
	xmlcontent, _ := xml.MarshalIndent(result, "", "    ")
	xmlbytes := append(xmlhead, xmlcontent...)

	//c.XML(http.StatusOK, result)
	c.String(http.StatusOK, string(xmlbytes))
	return
}

//获取注册码信息
func PhoneCodeInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}
	var phonecode PhoneCode
	if phonecode, err = service.GetPhoneCode(int32(id)); err != nil {
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(phonecode, "请求成功"))
}
