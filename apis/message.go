package apis

import (
	"photoshare/middleware"
	"photoshare/websocket"

	"github.com/gin-gonic/gin"
)

//route配置
func (router *Router) MessageRouteRegister() {
	router.GET("/message/ws", middleware.UserValidate, WSConn)
	router.OPTIONS("/message/ws", func(c *gin.Context) {})
	router.GET("/message/test", WSTest)
}

func WSConn(c *gin.Context) {
	user := GetUserInfo(c)
	websocket.StartClient(c.Writer, c.Request, &user)
}

func WSTest(c *gin.Context) {
	c.File("httptest/home.html")
}
