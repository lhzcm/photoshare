package main

import (
	"photoshare/config"
	db "photoshare/database"
	"strconv"
)

func main() {
	defer db.SqlDB.Close()
	router := initRouter()
	router.Run(":" + strconv.Itoa(config.Configs.Server.Port))
}
