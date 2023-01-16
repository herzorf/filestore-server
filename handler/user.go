package handler

import (
	"fmt"
	"github.com/herzorf/filestroe-server/db"
	"github.com/herzorf/filestroe-server/util"
	"net/http"
	"os"
)

const pwdSalt = "*#890"

func SignUpHandler(write http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		file, err := os.ReadFile("./static/view/signup.html")
		if err != nil {
			fmt.Println("文件读取错误", err)
			write.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = write.Write(file)
		if err != nil {
			fmt.Println("write err", err)
			write.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if request.Method == "POST" {
		err := request.ParseForm()
		if err != nil {
			fmt.Println("request.ParseForm err", err)
			write.WriteHeader(http.StatusInternalServerError)
			return
		}
		username := request.Form.Get("username")
		password := request.Form.Get("password")
		if len(password) < 5 || len(username) < 3 {
			_, err := write.Write([]byte("用户名，密码格式不对"))
			if err != nil {
				fmt.Println("write err", err)
				write.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		encPassword := util.Sha1([]byte(password + pwdSalt))
		suc := db.UserSignUp(username, encPassword)
		if !suc {
			write.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			_, err := write.Write([]byte("success sign in "))
			if err != nil {
				fmt.Println("write err ", err)
			}
		}
	}
}

func SignInHandler(write http.ResponseWriter, request *http.Request) {

}
