package db

import (
	"fmt"
	"github.com/herzorf/filestroe-server/db/mysql"
)

// UserSignUp 向数据库中添加用户名和密码
func UserSignUp(username string, password string) bool {
	prepare, err := mysql.ConnectDB().Prepare("INSERT INTO user (user_name,user_pwd) VALUES (?,?)")
	if err != nil {
		fmt.Println("prepare err", err)
		return false
	}
	defer func() {
		err2 := prepare.Close()
		if err2 != nil {
			fmt.Println("prepare close err", err)
		}
	}()
	exec, err := prepare.Exec(username, password)
	if err != nil {
		fmt.Println("insert err", err)
	}
	if affected, err := exec.RowsAffected(); err == nil && affected > 0 {
		return true
	}
	return false
}

func UserSignIn(username string, password string) bool {
	stmt, err := mysql.ConnectDB().Prepare("SELECT * from user WHERE user_name = ? LIMIT = 1")
	if err != nil {
		fmt.Println("prepare err", err)
		return false
	}
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println("query err", err)
		return false
	} else if rows == nil {
		fmt.Printf("username : %s not found\n", username)
		return false
	}
	fmt.Printf("%+v\n", rows)
	parseRows := mysql.ParseRows(rows)
	if len(parseRows) > 0 && string(parseRows[0]["user_pwd"].([]byte)) == password {
		return true
	}
	return false

}
