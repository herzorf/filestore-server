package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	DB, _ = sql.Open("mysql", "root:1234567@tcp(127.0.0.1:3307)/filestore")
	DB.SetMaxIdleConns(100)
	err := DB.Ping()
	if err != nil {
		panic(err)
	}
}

func ConnectDB() *sql.DB {
	fmt.Println(DB)
	return DB
}
