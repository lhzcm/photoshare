package main

import (
	. "photoshare/apis"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	router := gin.Default()
	//router.UserRouteRegister()
	RouteRegister(router)
	return router
}
