package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/herzorf/filestroe-server/db"
	"github.com/herzorf/filestroe-server/util"
	"log"
	"net/http"
	"os"
	"time"
)

const pwdSalt = "*#890"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUpHandler(c *gin.Context) {
	var user User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Fatal("bind err", err)
	}
	if len(user.Password) < 6 || len(user.Username) < 6 {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "用户名密码至少六位。",
			"data":    nil,
		})
	}
	encPassword := util.Sha1([]byte(user.Password + pwdSalt))
	err = db.UserSignUp(user.Username, encPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 0})
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

func IsTokenValid(username string, token string) bool {
	userToken, err := db.GetUserToken(username)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if userToken != token {
		return false
	}
	return true
}

func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
