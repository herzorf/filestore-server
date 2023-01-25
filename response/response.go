package response

import (
	"github.com/gin-gonic/gin"
)

func Response(ctx *gin.Context, httpCode int, code int, message string, data interface{}) {
	ctx.JSON(httpCode, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func Success(ctx *gin.Context, message string, data interface{}) {
	Response(ctx, 200, 0, message, data)
}

func Fail(ctx *gin.Context, message string, data interface{}) {
	Response(ctx, 200, -1, message, data)
}
