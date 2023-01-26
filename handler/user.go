package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/herzorf/filestroe-server/db"
	"github.com/herzorf/filestroe-server/response"
	"github.com/herzorf/filestroe-server/util"
	"log"
	"net/http"
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
		response.Fail(c, "用户名密码至少六位。", nil)
	}
	encPassword := util.Sha1([]byte(user.Password + pwdSalt))
	err = db.UserSignUp(user.Username, encPassword)
	if err != nil {
		response.Fail(c, "注册失败", nil)
	} else {
		response.Success(c, "注册成功", nil)
	}
}

func SignInHandler(c *gin.Context) {
	var user User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println("bind err", err)
		response.Response(c, http.StatusInternalServerError, -1, "", nil)
	}
	encPassword := util.Sha1([]byte(user.Password + pwdSalt))

	passwordCheck := db.UserSignIn(user.Username, encPassword)
	if !passwordCheck {
		response.Fail(c, "用户名密码错误", nil)
		return
	}
	token := GenToken(user.Username)
	updateToken := db.UpdateToken(user.Username, token)
	if !updateToken {
		response.Response(c, http.StatusInternalServerError, -1, "", nil)
		return
	} else {
		resp := struct {
			Location string
			Username string
			Token    string
		}{
			Location: c.ClientIP() + "/static/view/home.html",
			Username: user.Username,
			Token:    token,
		}
		response.Success(c, "登录成功", resp)
	}
}

func UserInfoHandler(c *gin.Context) {
	var user User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println("bind err ", err)
	}
	info, err := db.GetUserInfo(user.Username)
	if err != nil {
		response.Response(c, http.StatusForbidden, -1, "获取用户信息失败", nil)
		return
	}
	response.Success(c, "请求成功", info)
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
