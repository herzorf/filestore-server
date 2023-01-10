package handler

import (
	"fmt"
	"io"
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
		_, err = writer.Write(file)
		if err != nil {
			fmt.Println("文件写入错误", err)
		}
	} else if request.Method == "POST" {
		file, header, err := request.FormFile("file")
		if err != nil {
			fmt.Println("文件读取错误", err)
			return
		}
		defer func() {
			err := file.Close()
			if err != nil {
				fmt.Println("读取文件关闭错误", err)
			}
		}()
		newFile, err := os.Create("temp/" + header.Filename)
		if err != nil {
			fmt.Println("文件创建错误", err)
			return
		}
		defer func() {
			err := newFile.Close()
			if err != nil {
				fmt.Println("创建的文件关闭错误", err)
			}
		}()
		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Println("文件拷贝错误", err)
		}
		http.Redirect(writer, request, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSucHandler Upload Success
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "upload finished!")
	if err != nil {
		fmt.Println("io write err", err)
	}
}
