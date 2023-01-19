package main

import (
	"fmt"
	"github.com/herzorf/filestroe-server/cache/redis"
	"github.com/herzorf/filestroe-server/db/mysql"
	"github.com/herzorf/filestroe-server/handler"
	"net/http"
)

func main() {
	mysql.ConnectDB()
	pool := redis.RedisPool()
	fmt.Printf("%+v", pool)
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/file/fastUpload", handler.HTTPIntercepter(handler.TryFastUploadHandler))

	http.HandleFunc("/user/signup", handler.SignUpHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/filemeta", handler.UserFileQueryHandler)
	http.HandleFunc("/user/info", handler.HTTPIntercepter(handler.UserInfoHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("server start err %s\n", err.Error())
	}
}
