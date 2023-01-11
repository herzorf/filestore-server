package handler

import (
	"encoding/json"
	"fmt"
	"github.com/herzorf/filestroe-server/meta"
	"github.com/herzorf/filestroe-server/util"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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
		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: "temp/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-11 15:4:5"),
		}
		newFile, err := os.Create(fileMeta.Location)
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
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Println("文件拷贝错误", err)
		}
		_, err = newFile.Seek(0, 0)
		if err != nil {
			fmt.Println("file seek 错误", err)
		}
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)
		fmt.Printf("%+v\n", meta.FileMetas)
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

func GetFileMetaHandler(write http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		panic(err)
	}
	fileHash := request.Form["filehash"][0]
	fileMeta := meta.GetFileMeta(fileHash)
	marshal, err := json.Marshal(fileMeta)
	if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	_, err = write.Write(marshal)
	if err != nil {
		panic(err)
	}
}

func DownloadHandler(write http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		panic(err)
	}
	filesha1 := request.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(filesha1)
	fmt.Println(fileMeta.Location)
	file, err := os.Open(fileMeta.Location)
	if err != nil {
		log.Printf("文件打开错误%s", err)
		write.WriteHeader(http.StatusInternalServerError)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("文件关闭错误%s", err)
		}
	}()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("读取文件出错\n")
		write.WriteHeader(http.StatusInternalServerError)
	}
	_, err = write.Write(data)
	if err != nil {
		log.Printf("文件返回出错")
	}
}
