package handler

import (
	"fmt"
	"github.com/herzorf/filestroe-server/db"
	"github.com/herzorf/filestroe-server/util"
	"net/http"
	"os"
	"time"
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
	if request.Method == "GET" {
		file, err := os.ReadFile("./static/view/signin.html")
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
	} else {
		err := request.ParseForm()
		if err != nil {
			fmt.Println(err)
			write.WriteHeader(http.StatusInternalServerError)
		}
		username := request.Form.Get("username")
		password := request.Form.Get("password")

		encPassword := util.Sha1([]byte(password + pwdSalt))

		passwordCheck := db.UserSignIn(username, encPassword)
		if !passwordCheck {
			_, err := write.Write([]byte("用户名密码错误"))
			if err != nil {
				write.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		token := GenToken(username)
		updateToken := db.UpdateToken(username, token)
		if !updateToken {
			write.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			resp := util.RespMsg{
				Code: 0,
				Msg:  "ok",
				Data: struct {
					Location string
					Username string
					Token    string
				}{
					Location: request.Host + "/static/view/home.html",
					Username: username,
					Token:    token,
				},
			}
			_, err := write.Write(resp.JSONBytes())
			if err != nil {
				panic(err)
			}
		}
	}
}

func UserInfoHandler(write http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	username := request.Form.Get("username")
	info, err := db.GetUserInfo(username)
	if err != nil {
		write.WriteHeader(http.StatusForbidden)
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "ok",
		Data: info,
	}
	_, err = write.Write(resp.JSONBytes())
	if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
	}
}

func IsTokenValid(token string) bool {
	fmt.Println(token)
	return true
}

func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
