package database

import (
	"fmt"
	"log"
	"photoshare/config"
	. "photoshare/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var GormDB *gorm.DB

func init() {
	var conn config.MssqlConfig = config.Configs.Mssql
	connStr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		conn.Server, conn.User, conn.Password, conn.Port, conn.Database)
	db, err := gorm.Open("mssql", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	GormDB = db
	GormDB.SingularTable(true)
	//自动迁移
	GormDB.AutoMigrate(&PhoneCode{}, &Photo{}, &Publish{})
}
