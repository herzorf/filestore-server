package main

import (
	"fmt"
	"github.com/herzorf/filestroe-server/db/mysql"
	"github.com/herzorf/filestroe-server/handler"
	"net/http"
)

func main() {
	mysql.ConnectDB()
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/user/signup", handler.SignUpHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.UserInfoHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("server start err %s\n", err.Error())
	}
}
