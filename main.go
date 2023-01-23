package main

import (
	"fmt"
	"github.com/herzorf/filestroe-server/config/cos"
	"github.com/herzorf/filestroe-server/handler"
	"net/http"
)

func main() {
	//mysql.ConnectDB()
	//pool := redis.RedisPool()
	//fmt.Printf("%+v", pool)
	cos.GetUploadObjectUrl("萌兔兔.png")
	cos.ConnectCos()
	http.HandleFunc("/api/file/upload", handler.UploadHandler) //do
	//http.HandleFunc("/api/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/api/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/api/file/download", handler.DownloadHandler)
	http.HandleFunc("/api/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/api/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/api/file/fastUpload", handler.HTTPIntercepter(handler.TryFastUploadHandler))

	http.HandleFunc("/api/user/signup", handler.SignUpHandler)                          //do
	http.HandleFunc("/api/user/signin", handler.SignInHandler)                          //do
	http.HandleFunc("/api/user/filemeta", handler.UserFileQueryHandler)                 //do
	http.HandleFunc("/api/user/info", handler.HTTPIntercepter(handler.UserInfoHandler)) //do

	http.HandleFunc("/api/file/mpupload/init", handler.HTTPIntercepter(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/api/file/mpupload/uppart", handler.HTTPIntercepter(handler.UploadPartHandler))
	http.HandleFunc("/api/file/mpupload/complete", handler.HTTPIntercepter(handler.CompleteUploadHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("server start err %s\n", err.Error())
	}
}
