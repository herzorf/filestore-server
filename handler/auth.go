package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/herzorf/filestroe-server/response"
	"io"
	"log"
	"net/http"
)

type UserInfo struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

// HTTPIntercepter 请求拦截器
func HTTPIntercepter() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userInfo UserInfo
		requestBody, err := c.GetRawData()
		if err != nil {
			log.Println(err.Error())
		}
		err = json.Unmarshal(requestBody, &userInfo)
		if err != nil {
			log.Println("json Unmarshal err", err)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody)) // 关键点
		if len(userInfo.Username) < 6 || !IsTokenValid(userInfo.Username, userInfo.Token) {
			response.Response(c, http.StatusForbidden, 403, "token错误或用户名格式不对", nil)
			return
		}
		c.Next()
	}
}
