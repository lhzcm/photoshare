package main

import (
	"photoshare/config"
	db "photoshare/database"
	"photoshare/websocket"
	"strconv"
)

func main() {
	defer db.SqlDB.Close()
	router := initRouter()
	go websocket.Start()
	router.Run(":" + strconv.Itoa(config.Configs.Server.Port))
}
