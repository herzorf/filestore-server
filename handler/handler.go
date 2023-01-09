package handler

import (
	"fmt"
	"net/http"
	"os"
)

// UploadHandler 处理上传文件
func UploadHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		file, err := os.ReadFile("./static/view/index.html")
		if err != nil {
			fmt.Printf("文件读取错误 %s", err.Error())
		}
		write, err := writer.Write(file)
		if err != nil {
			fmt.Println("文件写入错误", err)
		}
		fmt.Println(write)
	} else if request.Method == "POST" {
		fmt.Println(request.Method)

	}
}
