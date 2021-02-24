package database

import (
	"fmt"
	"log"
	"photoshare/config"
	. "photoshare/models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func init() {
	var conn config.MssqlConfig = config.Configs.Mssql
	connStr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		conn.Server, conn.User, conn.Password, conn.Port, conn.Database)
	dbconfig := sqlserver.Config{
		DriverName: "sqlserver",
		DSN:        connStr,
	}
	db, err := gorm.Open(sqlserver.New(dbconfig), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	GormDB = db
	//GormDB.SingularTable(true)
	//自动迁移
	GormDB.AutoMigrate(&PhoneCode{}, &Photo{}, &Publish{}, &Comment{}, &Praise{}, &Invitation{}, &Friend{}, &Message{})
}
