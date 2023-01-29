package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "xxx",
	})
}
