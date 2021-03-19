package apis

import (
	"net/http"
	"photoshare/config"

	"github.com/gin-gonic/gin"
)

type Router gin.Engine

func RouteRegister(e *gin.Engine) {
	router := (*Router)(e)
	//允许跨域
	router.Use(CORS)
	//允许跨域option
	router.OPTIONS("/*param", func(c *gin.Context) {})
	router.UserRouteRegister()
	router.PublishRouteRegister()
	router.InvitationRouteRegister()
	router.FriendRouteRegister()
	router.MessageRouteRegister()
	//router.HomeRouteRegister()

	//静态文件
	router.StaticFS("/static/image", http.Dir(config.Configs.Static.PublishImgPath))
}

//跨域设置
func CORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", config.Configs.Server.Corshost)
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, Origin, X-Requested-With, Accept")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTION")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Vary", "Origin")
	// c.Header("Access-Control-Max-Age", "3600")
}
