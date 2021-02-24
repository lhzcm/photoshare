package apis

import (
	"net/http"
	"photoshare/middleware"
	. "photoshare/models"
	"photoshare/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

//route配置
func (router *Router) FriendRouteRegister() {
	router.GET("/friend", middleware.UserValidate, FriendsList)
	router.DELETE("/friend/:friendid", middleware.UserValidate, DelFriend)
}

//获取好友列表
func FriendsList(c *gin.Context) {
	user := GetUserInfo(c)
	friends, err := service.FriendsList(user.Id)

	if err != nil {
		c.JSON(http.StatusOK, Fail("获取好友列表失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(friends, "请求成功"))
}

//删除好友
func DelFriend(c *gin.Context) {
	var friendid int
	var err error

	user := GetUserInfo(c)

	if friendid, err = strconv.Atoi(c.Param("friendid")); err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}
	if err = service.DelFriend(user.Id, int32(friendid)); err != nil {
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(nil, "删除成功"))
}
