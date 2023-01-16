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
	return DB
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println(err)
			return records
		}

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}
