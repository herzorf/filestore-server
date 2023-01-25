package route

import (
	"github.com/gin-gonic/gin"
	"github.com/herzorf/filestroe-server/handler"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/api/user/signup", handler.SignUpHandler)

	return r
}
