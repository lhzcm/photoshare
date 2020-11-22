package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

var SqlDB *sql.DB

const (
	server   = "112.74.48.21"
	port     = 63354
	user     = "sa"
	password = "Lkl888888"
	database = "photoshare"
)

func init() {
	var err error
	connStr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
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
