package apis

import (
	"log"
	"net/http"
	"photoshare/config"
	"photoshare/middleware"
	"photoshare/service"
	"photoshare/utility"
	"photoshare/websocket"
	"strconv"

	. "photoshare/models"

	"github.com/gin-gonic/gin"
)

//route配置
func (router *Router) MessageRouteRegister() {
	router.GET("/message/ws", middleware.UserValidate, WSConn)
	router.GET("/messages/:receiverid/:id", middleware.UserValidate, GetMessageList)
	router.POST("/messages/uploadimg", middleware.UserValidate, UploadImg)
}

//websocket 连接
func WSConn(c *gin.Context) {
	user := GetUserInfo(c)
	websocket.StartClient(c.Writer, c.Request, &user)
}

//获取消息列表
func GetMessageList(c *gin.Context) {
	senderidstr, hasreceiverid := c.Params.Get("receiverid")
	if !hasreceiverid {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}
	senderid, rerr := strconv.Atoi(senderidstr)
	if rerr != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}

	idstr, hasid := c.Params.Get("id")
	if !hasid {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}

	user := GetUserInfo(c)
	var messages []Message
	if messages, err = service.GetMessageList(int(user.Id), senderid, id); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, Fail("查询出错"))
		return
	}
	c.JSON(http.StatusOK, Success(messages, "查询成功"))
}

//图片上传
func UploadImg(c *gin.Context) {
	file, err := c.FormFile("img")
	if err != nil {
		c.JSON(http.StatusOK, Fail("文件上传失败"))
		return
	}
	if !utility.ValidateImgName(file.Filename) {
		c.JSON(http.StatusOK, Fail("不支持的图片类型"))
		return
	}
	if file.Size > config.Configs.Static.MaxUploadImgSize {
		c.JSON(http.StatusOK, Fail("上传失败，上传文件过大"))
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, Fail("上传失败"))
		return
	}
	fileName := utility.GetGUID().Hex() + "." + utility.GetFileExtend(file.Filename)
	if err = c.SaveUploadedFile(file, "./images/uploadimg/"+fileName); err != nil {
		c.JSON(http.StatusOK, Fail("上传失败"))
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, Success(fileName, "请求成功"))
}
