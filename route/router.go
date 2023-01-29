package route

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/herzorf/filestroe-server/handler"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var FS embed.FS

func Router() *gin.Engine {
	r := gin.Default()
	staticFiles, _ := fs.Sub(FS, "dist")
	r.StaticFS("/", http.FS(staticFiles))
	//r.GET("/signin", handler.GetPage)
	r.POST("/api/user/signup", handler.SignUpHandler)
	r.POST("/api/user/signin", handler.SignInHandler)
	r.POST("/api/user/info", handler.UserInfoHandler)

	r.Use(handler.HTTPIntercepter())
	r.POST("/api/user/filemeta", handler.UserFileQueryHandler)
	r.POST("/api/file/upload", handler.UploadHandler)
	r.POST("/api/file/delete", handler.FileDeleteHandler)

	return r
}
