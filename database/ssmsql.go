package database

import (
	"context"
	"database/sql"
	"fmt"
	"photoshare/config"

	_ "github.com/denisenkom/go-mssqldb"
)

var SqlDB *sql.DB

func init() {
	var err error
	var conn config.MssqlConfig = config.Configs.Mssql
	connStr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		conn.Server, conn.User, conn.Password, conn.Port, conn.Database)
	SqlDB, err = sql.Open("sqlserver", connStr)
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	err = SqlDB.PingContext(ctx)
	if err != nil {
		panic(err.Error())
	}
}
