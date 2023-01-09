package handler

import (
	"fmt"
	"net/http"
)

// UploadHandler 处理上传文件
func UploadHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.Method)
}
