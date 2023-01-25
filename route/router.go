package route

import (
	"github.com/gin-gonic/gin"
	"github.com/herzorf/filestroe-server/handler"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/api/user/signup", handler.SignUpHandler)
	r.POST("/api/user/signin", handler.SignInHandler)
	r.POST("/api/user/info", handler.UserInfoHandler)

	return r
}
