package main

import (
	"github.com/gin-gonic/gin"
	. "lkl.photoshare/apis"
)

func initRouter() *gin.Engine {
	router := gin.Default()
	//router.UserRouteRegister()
	RouteRegister(router)
	return router
}
