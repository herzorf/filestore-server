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
