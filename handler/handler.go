package handler

import (
	"encoding/json"
	"fmt"
	"github.com/herzorf/filestroe-server/config/cos"
	"github.com/herzorf/filestroe-server/db"
	"github.com/herzorf/filestroe-server/meta"
	"github.com/herzorf/filestroe-server/util"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// UploadHandler 处理上传文件
func UploadHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
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
		meta.UpdateFileMetaDB(fileMeta)
		fd, err := os.Open(fileMeta.Location)
		if err != nil {
			panic(err)
		}
		defer fd.Close()
		err = cos.PutFileObject(fd, fileMeta.FileName)
		if err != nil {
			fmt.Println("put object err", err)
		}
		_ = request.ParseForm()
		username := request.Form.Get("username")
		finished := db.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if !finished {
			_, err := writer.Write([]byte("Upload fail"))
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
			}
			return
		} else {
			_, err = writer.Write(util.NewRespMsg(0, "success upload", nil).JSONBytes())
		}
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
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
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

func UserFileQueryHandler(write http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
		return
	}
	username := request.Form.Get("username")
	limit, _ := strconv.Atoi(request.Form.Get("limit"))
	metas, err := db.QueryUserFileMetas(username, limit)
	//marshal, _ := json.Marshal(metas)
	if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		_, err = write.Write(util.NewRespMsg(0, "ok", metas).JSONBytes())
		if err != nil {
			write.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
func DownloadHandler(write http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		panic(err)
	}
	filesha1 := request.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(filesha1)
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

func FileMetaUpdateHandler(write http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		panic(err)
	}
	opType := request.Form.Get("op")
	fileSha1 := request.Form.Get("filehash")
	newFileName := request.Form.Get("filename")

	if opType != "0" {
		write.WriteHeader(http.StatusForbidden)
		return
	}

	if request.Method != "POST" {
		write.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fileMeta := meta.GetFileMeta(fileSha1)
	fileMeta.FileName = newFileName
	meta.UpdateFileMetaDB(fileMeta)

	marshal, err := json.Marshal(fileMeta)
	if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
		return
	}
	write.WriteHeader(http.StatusOK)
	_, err = write.Write(marshal)
	if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func FileDeleteHandler(write http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		panic(err)
	}
	fileSha1 := request.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fileSha1)
	err = os.Remove(fileMeta.Location)
	if err != nil {
		fmt.Println("文件删除错误", err)
	}
	meta.RemoveFileMeta(fileSha1)
	write.WriteHeader(http.StatusOK)
}

func TryFastUploadHandler(write http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Println(err)
		write.WriteHeader(http.StatusInternalServerError)
		return
	}
	username := request.Form.Get("username")
	fileHash := request.Form.Get("fileHash")
	fileName := request.Form.Get("fileName")
	fileSize, _ := strconv.Atoi(request.Form.Get("fileSize"))
	fileMeta, err := db.OnGetFileMeta(fileHash)
	if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
		return
	}
	if fileMeta == nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "妙传失败，请使用普通上传接口",
			Data: nil,
		}
		_, err = write.Write(resp.JSONBytes())
		return
	}
	finished := db.OnUserFileUploadFinished(username, fileHash, fileName, int64(fileSize))
	if finished {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "妙传成功",
			Data: nil,
		}
		_, err = write.Write(resp.JSONBytes())
		return
	} else {
		resp := util.RespMsg{
			Code: -2,
			Msg:  "妙传失败",
			Data: nil,
		}
		_, err = write.Write(resp.JSONBytes())
		return
	}
}
