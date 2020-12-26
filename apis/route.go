package apis

import "github.com/gin-gonic/gin"

type Router gin.Engine

func RouteRegister(e *gin.Engine) {
	router := (*Router)(e)
	router.UserRouteRegister()
	router.PublishRouteRegister()
}
