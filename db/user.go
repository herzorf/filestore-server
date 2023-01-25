package db

import (
	"fmt"
	"github.com/herzorf/filestroe-server/db/mysql"
	"log"
)

// UserSignUp 向数据库中添加用户名和密码
func UserSignUp(username string, password string) error {
	prepare, err := mysql.ConnectDB().Prepare("INSERT INTO user (user_name,user_pwd) VALUES (?,?)")
	if err != nil {
		log.Println("prepare err", err)
		return err
	}
	defer func() {
		err2 := prepare.Close()
		if err2 != nil {
			log.Println("prepare close err", err)
		}
	}()
	exec, err := prepare.Exec(username, password)
	if err != nil {
		log.Println("insert err", err)
		return err
	}
	if affected, err := exec.RowsAffected(); err == nil && affected > 0 {
		return nil
	}
	return err
}

func UserSignIn(username string, password string) bool {
	stmt, err := mysql.ConnectDB().Prepare("SELECT * from user WHERE user_name = ? LIMIT 1")
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
	parseRows := mysql.ParseRows(rows)
	if len(parseRows) > 0 && string(parseRows[0]["user_pwd"].([]byte)) == password {
		return true
	}
	return false
}

func UpdateToken(username string, token string) bool {
	stmt, err := mysql.ConnectDB().Prepare("replace into user_token (user_name,user_token) values (?,?)")
	if err != nil {
		fmt.Println("prepare err", err)
		return false
	}
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println("replace err", err)
		return false
	}
	defer func() {
		err2 := stmt.Close()
		if err2 != nil {
			panic(err2)
		}
	}()
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignUpAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	prepare, err := mysql.ConnectDB().Prepare("SELECT user_name,signup_at FROM  user WHERE user_name=? LIMIT 1")
	if err != nil {
		fmt.Println("prepare err", err)
		return user, err
	}
	if err != nil {
		fmt.Println("prepare err", err)
		return user, nil
	}
	defer func() {
		err2 := prepare.Close()
		if err2 != nil {
			panic(err2)
		}
	}()
	err = prepare.QueryRow(username).Scan(&user.Username, &user.SignUpAt)
	if err != nil {
		return user, err
	}
	return user, err
}

func GetUserToken(username string) (string, error) {
	token := ""
	prepare, err := mysql.ConnectDB().Prepare("SELECT user_token FROM user_token WHERE user_name =? LIMIT 1")
	if err != nil {
		fmt.Println("prepare err", err)
		return token, err
	}
	defer func() {
		err2 := prepare.Close()
		if err2 != nil {
			panic(err2)
		}
	}()
	err = prepare.QueryRow(username).Scan(&token)
	if err != nil {
		return token, err
	}
	return token, nil
}
