package main

import (
	"fmt"
	"github.com/herzorf/filestroe-server/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("server start err %s\n", err.Error())
	}
}
