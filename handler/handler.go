package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/herzorf/filestroe-server/config/cos"
	"github.com/herzorf/filestroe-server/db"
	"github.com/herzorf/filestroe-server/meta"
	"github.com/herzorf/filestroe-server/response"
	"github.com/herzorf/filestroe-server/util"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadHandler(c *gin.Context) {
	file, _ := c.FormFile("file")
	fileMeta := meta.FileMeta{
		FileName: file.Filename,
		Location: "temp/" + file.Filename,
		FileSize: file.Size,
		UploadAt: time.Now().Format("2006-01-11 15:4:5"),
	}
	openFile, err := file.Open()
	body := &bytes.Buffer{}
	_, err = io.Copy(body, openFile)
	if err != nil {
		log.Println("copy err", err)
	}
	reader := bytes.NewReader(body.Bytes())
	fileMeta.FileSha1 = util.FileSha1(reader)
	meta.UpdateFileMetaDB(fileMeta)
	contentType := file.Header.Get("Content-Type")
	err = cos.PutFileObject(body, fileMeta.FileSha1, contentType)
	if err != nil {
		fmt.Println("put object err", err)
	}
	username := c.PostForm("username")
	finished := db.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
	if !finished {
		response.Fail(c, "上传失败", nil)
	} else {
		response.Success(c, "上传成功", nil)
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

func UserFileQueryHandler(c *gin.Context) {
	type QueryFileMeta struct {
		Username string `json:"username"`
		Limit    int    `json:"limit"`
	}
	var queryFileMeta QueryFileMeta
	err := c.ShouldBindJSON(&queryFileMeta)
	if err != nil {
		log.Println("gin bind err111", err)
	}

	metas, err := db.QueryUserFileMetas(queryFileMeta.Username, queryFileMeta.Limit)
	if err != nil {
		response.Fail(c, "查询失败", nil)
		return
	} else {
		response.Success(c, "请求成功", metas)
		return
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
