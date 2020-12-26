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

	"github.com/gin-gonic/gin"
)

//route配置
func (router *Router) PublishRouteRegister() {
	router.POST("/publish/uploadimg", middleware.UserValidate, ImageUpload)
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
