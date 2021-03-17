package apis

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
	"strings"

	"photoshare/middleware"
	. "photoshare/models"
	"photoshare/service"
	"photoshare/utility"

	"github.com/gin-gonic/gin"
)

//route配置
func (router *Router) UserRouteRegister() {
	router.Use(gin.Logger())
	router.GET("/user/code/:phone/:type", GetPhoneCode)
	// router.GET("/user/code/:id", PhoneCodeInfo)
	router.GET("/user/callback", CallBackPhoneCode)
	router.POST("/user", UserRegister)
	router.POST("/user/login", UserLogin)
	router.GET("/user/info", middleware.UserValidate, UserInfo)
	router.GET("/user/login/:phone/:code", UserLoginByCode)
	router.GET("/user/register/:phone/:code", UserRegisterByCode)
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

	if phoneCode, err := service.ValidatePhoneCode(user.Phone, user.Code, 1); err != nil || phoneCode.Status != 1 {
		if err != nil {
			c.JSON(http.StatusOK, Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Fail("注册失败，手机号码没有验证"))
		return
	}

	user.Password = utility.EncryptPassword(user.Password)

	if err := service.UserRegister(&user); err != nil {
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	if _, token, err := service.UserLogin(user.Id, "", user.Password); err == nil {
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
	if user, token, err = service.UserLogin(user.Id, user.Phone, user.Password); err != nil {
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	c.SetCookie("token", token, 60*24*30, "/", c.Request.URL.Host, false, false)
	c.JSON(http.StatusOK, Success(user, "登录成功"))
}

//用户验证码登录
func UserLoginByCode(c *gin.Context) {
	var (
		err       error
		user      User
		code      int
		phoneCode PhoneCode
	)
	data := make(map[string]interface{}, 2)

	if code, err = strconv.Atoi(c.Param("code")); err != nil {
		c.JSON(http.StatusOK, Fail("验证码有误"))
		return
	}
	user.Phone = c.Param("phone")
	if phoneCode, err = service.ValidatePhoneCode(user.Phone, code, 0); err != nil || phoneCode.Id <= 0 {
		c.JSON(http.StatusOK, Fail("获取验证码有误"))
		return
	}
	if phoneCode.Status == 0 {
		data["Code"] = phoneCode
		data["User"] = nil
		c.JSON(http.StatusOK, Success(data, "请求成"))
		return
	}
	if err = service.GetUserInfoByPhone(&user); err != nil {
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}

	data["Code"] = phoneCode
	data["User"] = user
	c.SetCookie("token", user.Token, 60*24*30, "/", c.Request.URL.Host, false, false)
	c.JSON(http.StatusOK, Success(data, "登录成功"))
}

//用户注册判断验证是否成功
func UserRegisterByCode(c *gin.Context) {

	var phoneCode PhoneCode
	var err error
	var code int

	if code, err = strconv.Atoi(c.Param("code")); err != nil {
		c.JSON(http.StatusOK, Fail("验证码有误"))
		return
	}
	phoneCode.Phone = c.Param("phone")

	if phoneCode, err = service.ValidatePhoneCode(phoneCode.Phone, code, 1); err != nil || phoneCode.Id <= 0 {
		if err != nil {
			c.JSON(http.StatusOK, Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Fail("获取验证码有误"))
		return
	}

	c.JSON(http.StatusOK, Success(phoneCode, "请求成功"))
	return
}

//用户登录/注册返回code
func GetPhoneCode(c *gin.Context) {
	phone := c.Param("phone")
	if len(phone) != 11 {
		c.JSON(http.StatusOK, Fail("手机号码有误"))
		return
	}
	types := c.Param("type")
	var typeint int
	if types == "0" {
		typeint = 0
	} else if types == "1" {
		typeint = 1
	} else {
		c.JSON(http.StatusOK, Fail("验证码类型有误"))
		return
	}

	if code, err := service.CreatePhoneCode(phone, typeint); err == nil {
		c.JSON(http.StatusOK, Success(code, "成功"))
		return
	} else {
		c.JSON(http.StatusOK, Fail(err.Error()))
	}
}

//第三方回调Phone注册码
func CallBackPhoneCode(c *gin.Context) {
	xmlHead := "<?xml version=\"1.0\" encoding=\"GB2312\"?>\n"
	xmlhead := []byte(xmlHead)
	phone := c.Query("mobile")
	result := KeliResult{Version: "1.0"}
	strs := strings.Split(c.Query("content"), "xin#")

	if len(strs) != 2 {
		result.ResultInfo = "param error"
		xmlcontent, _ := xml.MarshalIndent(result, "", "    ")
		xmlbytes := append(xmlhead, xmlcontent...)
		c.String(http.StatusOK, string(xmlbytes))
		return
	}

	code, err := strconv.ParseInt(strs[1], 10, 32)
	if err != nil {
		result.ResultInfo = "code error"
		xmlcontent, _ := xml.MarshalIndent(result, "", "    ")
		xmlbytes := append(xmlhead, xmlcontent...)
		c.String(http.StatusOK, string(xmlbytes))
		return
	}

	if len(phone) != 11 {
		result.ResultInfo = "phone error"
		xmlcontent, _ := xml.MarshalIndent(result, "", "    ")
		xmlbytes := append(xmlhead, xmlcontent...)
		c.String(http.StatusOK, string(xmlbytes))
		return
	}

	var rowcount int64
	if rowcount, err = service.UpdatePhoneCode(int32(code), phone); err != nil {
		result.ResultInfo = err.Error()
		xmlcontent, _ := xml.MarshalIndent(result, "", "    ")
		xmlbytes := append(xmlhead, xmlcontent...)
		c.String(http.StatusOK, string(xmlbytes))
		return
	}
	if rowcount <= 0 {
		result.ResultInfo = "update error"
		xmlcontent, _ := xml.MarshalIndent(result, "", "    ")
		xmlbytes := append(xmlhead, xmlcontent...)
		c.String(http.StatusOK, string(xmlbytes))
		return
	}

	result.Result = 1
	result.ResultInfo = "success"
	xmlcontent, _ := xml.MarshalIndent(result, "", "    ")
	xmlbytes := append(xmlhead, xmlcontent...)
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
