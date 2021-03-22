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
func (router *Router) InvitationRouteRegister() {
	router.POST("/invite", middleware.UserValidate, Invite)
	router.PATCH("/invite", middleware.UserValidate, InviteAccept)
	router.GET("/invite/:page/:pagesize/:status", middleware.UserValidate, InviteList)
}

//好友邀请
func Invite(c *gin.Context) {
	var invitation Invitation
	if err := c.BindJSON(&invitation); err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}
	user := GetUserInfo(c)
	invitation.Userid = user.Id
	invitation.Status = 0

	if err := service.Invite(&invitation); err != nil {
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success("邀请发送成功", "请求成功"))
	return
}

//接受好友邀请
func InviteAccept(c *gin.Context) {
	var invitation Invitation
	if err := c.BindJSON(&invitation); err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}

	if invitation.Status > 1 && invitation.Status < -1 {
		c.JSON(http.StatusOK, Fail("更新状态有误"))
		return
	}
	if err := service.InviteAccept(invitation.Id, GetUserInfo(c).Id, int16(invitation.Status)); err != nil {
		c.JSON(http.StatusOK, Fail(err.Error()))
		return
	}
	if invitation.Status == 1 {
		c.JSON(http.StatusOK, Success(nil, "接受成功"))
		return
	}
	c.JSON(http.StatusOK, Success(nil, "拒绝成功"))
}

//好友邀请列表
func InviteList(c *gin.Context) {
	var page, pagesize, status int
	var err error
	var invitations []Invitation
	var count int64

	if page, err = strconv.Atoi(c.Param("page")); err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}

	if pagesize, err = strconv.Atoi(c.Param("pagesize")); err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}

	if status, err = strconv.Atoi(c.Param("status")); err != nil {
		c.JSON(http.StatusOK, Fail("参数错误"))
		return
	}

	if invitations, count, err = service.InviteList(GetUserInfo(c).Id, int16(status), page, pagesize); err != nil {
		c.JSON(http.StatusOK, Fail("查询错误"))
		return
	}
	var result map[string]interface{} = make(map[string]interface{}, 2)
	result["list"] = invitations
	result["count"] = count
	c.JSON(http.StatusOK, Success(result, "查询成功"))
}
