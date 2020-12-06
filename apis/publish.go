package apis

import (
	"photoshare/middleware"

	"github.com/gin-gonic/gin"
)

//route配置
func (router *Router) PublishRouteRegister() {
	router.GET("/user/code", GetPhoneCode)
	router.GET("/user/code/:id", PhoneCodeInfo)
	router.GET("/user/callback", CallBackPhoneCode)
	router.POST("/user", UserRegister)
	router.POST("/user/login", UserLogin)
	router.GET("/user/info", middleware.UserValidate, UserInfo)
}

//图片上传接口
func ImageUpload(c *gin.Context) {

}
