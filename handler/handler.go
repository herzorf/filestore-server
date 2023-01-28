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
	for index, _ := range metas {
		url := cos.GetUploadObjectUrl(metas[index].FileHash)
		metas[index].Location = url.String()
	}
	if err != nil {
		response.Fail(c, "查询失败", nil)
		return
	} else {
		response.Success(c, "请求成功", metas)
		return
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

type DeleteFileRequest struct {
	Filehash string `json:"filehash"`
}

func FileDeleteHandler(c *gin.Context) {
	var deleteFileRequest DeleteFileRequest
	err := c.ShouldBindJSON(&deleteFileRequest)
	if err != nil {
		log.Println("gin bind err", err)
	}
	log.Println(deleteFileRequest.Filehash)
	err = cos.DeleteFileObject(deleteFileRequest.Filehash)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, -1, "删除失败", nil)
		return
	}
	response.Success(c, "删除成功", nil)
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
