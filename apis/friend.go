package apis

import (
	"net/http"
	"photoshare/middleware"
	. "photoshare/models"
	"photoshare/service"

	"github.com/gin-gonic/gin"
)

//route配置
func (router *Router) FriendRouteRegister() {
	router.GET("/friend", middleware.UserValidate, FriendsList)
}

func FriendsList(c *gin.Context) {
	user := GetUserInfo(c)
	friends, err := service.FriendsList(user.Id)

	if err != nil {
		c.JSON(http.StatusOK, Fail("获取好友列表失败"))
		return
	}
	c.JSON(http.StatusOK, Success(friends, "请求成功"))
}
