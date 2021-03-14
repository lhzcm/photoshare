package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//前端静态文件router
//route配置
func (router *Router) HomeRouteRegister() {
	router.GET("/home/*path", func(c *gin.Context) {
		c.HTML(http.StatusOK, "httptest/home.html", gin.H{})
	})
}
