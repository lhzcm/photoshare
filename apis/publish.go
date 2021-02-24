package apis

import (
	"log"
	"mime/multipart"
	"net/http"
	"photoshare/config"
	"photoshare/middleware"
	. "photoshare/models"
	"photoshare/service"
	"photoshare/utility"
	"strconv"

	"github.com/gin-gonic/gin"
)

//route配置
func (router *Router) PublishRouteRegister() {
	router.POST("/publish/uploadimg", middleware.UserValidate, ImageUpload)
	router.POST("/publish", middleware.UserValidate, Publishing)
	router.DELETE("/publish/:id", middleware.UserValidate, PublishDelete)
	router.GET("/publish/:page/:pagesize", middleware.UserValidate, GetPublishList)
	router.POST("/praise/:id/:type/:ispraise", middleware.UserValidate, PublishPraise)
	router.POST("/comment", middleware.UserValidate, PublishComment)
}

//图片上传接口
func ImageUpload(c *gin.Context) {
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
	var content multipart.File
	content, err = file.Open()
	if err != nil {
		c.JSON(http.StatusOK, Fail("上传失败"))
		return
	}
	exif := utility.GetImgExif(content)
	fileName := utility.GetGUID().Hex() + "." + utility.GetFileExtend(file.Filename)
	if err = c.SaveUploadedFile(file, "./images/uploadimg/"+fileName); err != nil {
		c.JSON(http.StatusOK, Fail("上传失败"))
		log.Println(err)
		return
	}

	photo := &Photo{
		Imgurl: fileName,
		Info:   exif,
		Userid: GetUserInfo(c).Id,
	}

	if err = service.AddUploadImg(photo); err != nil {
		c.JSON(http.StatusOK, Fail("上传失败"))
		return
	}
	c.JSON(http.StatusOK, Success(photo, "请求成功"))
}

//发布
func Publishing(c *gin.Context) {
	publish := Publish{}
	c.BindJSON(&publish)

	user := GetUserInfo(c)

	publish.Userid = user.Id

	if length := len(publish.Content); length < 5 || length > 512 {
		c.JSON(http.StatusOK, Fail("内容文字不能小于5个字且不能大于512个字"))
		return
	}
	imgids, err := utility.StringToIntArray(publish.Imgs, ",")
	if err != nil {
		c.JSON(http.StatusOK, Fail("上传的图片有误"))
		return
	}

	if len(imgids) > 0 && service.PhotoIsUser(imgids, int(user.Id)) < len(imgids) {
		c.JSON(http.StatusOK, Fail("上传图片有误，不是自己上传的图片"))
		return
	}

	err = service.SavePublish(&publish, imgids)
	if err != nil {
		c.JSON(http.StatusOK, Fail("动态发布失败"))
		return
	}

	c.JSON(http.StatusOK, Success(publish, "动态发布成功"))
}

//删除动态
func PublishDelete(c *gin.Context) {
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
	if err = service.DeletePublish(int32(id), user.Id); err != nil {
		c.JSON(http.StatusOK, Fail("删除失败"))
		return
	}
	c.JSON(http.StatusOK, Success(nil, "删除成功"))
}

//获取动态列表
func GetPublishList(c *gin.Context) {

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}

	var pagesize int
	pagesize, err = strconv.Atoi(c.Param("pagesize"))
	if err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}

	var result map[string]interface{} = make(map[string]interface{}, 2)
	user := GetUserInfo(c)

	publishs, total, err := service.GetPublishList(int(user.Id), page, pagesize)
	if err != nil {
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	result["publishs"] = publishs
	result["total"] = total

	c.JSON(http.StatusOK, Success(result, "请求成功"))
}

//用户点赞或者取消点赞
func PublishPraise(c *gin.Context) {
	var id, ptype, ispraise int
	var err error

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}
	if ptype, err = strconv.Atoi(c.Param("type")); err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}
	if ispraise, err = strconv.Atoi(c.Param("ispraise")); err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}
	user := GetUserInfo(c)

	if ispraise == 0 {
		if err = service.PublishUnPraise(user.Id, id, ptype); err != nil {
			c.JSON(http.StatusOK, Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Success(nil, "取消点赞成功"))
	} else {
		if err = service.PublishPraise(user.Id, id, ptype); err != nil {
			c.JSON(http.StatusOK, Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Success(nil, "点赞成功"))
	}
}

//动态评论
func PublishComment(c *gin.Context) {
	var comment Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}
	if len(comment.Content) <= 5 || len(comment.Content) > 256 {
		c.JSON(http.StatusOK, Fail("评论失败，评论字数不能小于5个字也不能大于256个字"))
		return
	}
	comment.Praise = 0
	comment.Comments = 0
	comment.Status = 0
	comment.Userid = GetUserInfo(c).Id

	if err := service.PublishComment(&comment); err != nil {
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(comment, "评论成功"))
}
