package handler

import (
	"fmt"
	"net/http"
)

// UploadHandler 处理上传文件
func UploadHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		fmt.Println(request.Method)
	} else if request.Method == "POST" {
		fmt.Println(request.Method)

	}
}
