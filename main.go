package main

import (
	db "lkl.photoshare/database"
)

func main() {
	defer db.SqlDB.Close()
	router := initRouter()
	router.Run(":8080")
}
