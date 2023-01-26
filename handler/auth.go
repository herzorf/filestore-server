package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/herzorf/filestroe-server/response"
	"net/http"
)

// HTTPIntercepter 请求拦截器
func HTTPIntercepter() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		token := c.GetHeader("token")
		if len(username) < 6 || !IsTokenValid(username, token) {
			c.Abort()
			response.Response(c, http.StatusForbidden, 403, username+token, nil)
			return
		} else {
			c.Next()
		}
	}
}
