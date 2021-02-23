package apis

import "github.com/gin-gonic/gin"

type Router gin.Engine

func RouteRegister(e *gin.Engine) {
	router := (*Router)(e)
	//允许跨域
	router.Use(CORS)
	router.UserRouteRegister()
	router.PublishRouteRegister()
	router.InvitationRouteRegister()
	router.FriendRouteRegister()
}

//跨域设置
func CORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

}
