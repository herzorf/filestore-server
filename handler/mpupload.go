package handler

import (
	"fmt"
	redis2 "github.com/gomodule/redigo/redis"
	"github.com/herzorf/filestroe-server/cache/redis"
	"github.com/herzorf/filestroe-server/db"
	"github.com/herzorf/filestroe-server/util"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type MultiPartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

func InitialMultipartUploadHandler(write http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Println("parseForm err", err)
		write.WriteHeader(http.StatusInternalServerError)
		return
	}
	username := request.Form.Get("username")
	fileHash := request.Form.Get("filehash")
	fileSize, err := strconv.Atoi(request.Form.Get("filesize"))
	if err != nil {
		_, err = write.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}
	rConn := redis.RedisPool().Get()
	defer func() {
		err2 := rConn.Close()
		if err2 != nil {
			panic(err2)
		}
	}()
	upinfo := MultiPartUploadInfo{
		FileHash:   fileHash,
		FileSize:   fileSize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, // 5M大小
		ChunkCount: int(math.Ceil(float64(fileSize) / (5 * 1024 * 1024))),
	}
	_, err = rConn.Do("HSET", "MP_"+upinfo.UploadID, "chunkcount", upinfo.ChunkCount)
	_, err = rConn.Do("HSET", "MP_"+upinfo.UploadID, "filehash", upinfo.FileHash)
	_, err = rConn.Do("HSET", "MP_"+upinfo.UploadID, "filesize", upinfo.FileSize)
	_, err = write.Write(util.NewRespMsg(0, "ok", nil).JSONBytes())
}

func UploadPartHandler(write http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Println("parseForm err", err)
		write.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = request.Form.Get("username")
	uploadID := request.Form.Get("uploadid")
	chunkIndex := request.Form.Get("index")

	rConn := redis.RedisPool().Get()

	defer func() {
		err2 := rConn.Close()
		if err2 != nil {
			panic(err2)
		}
	}()
	fpath := "/data/" + uploadID + "/" + chunkIndex
	err = os.MkdirAll(path.Dir(fpath), 0744)
	if err != nil {
		fmt.Println("MKdirAll err", err)
	}
	create, err := os.Create(fpath)
	if err != nil {
		fmt.Println("os create err", err)
		_, err = write.Write(util.NewRespMsg(-1, "upload part failed", nil).JSONBytes())
		return
	}
	defer func() {
		err2 := create.Close()
		if err2 != nil {
			panic(err2)
		}
	}()

	buf := make([]byte, 1024*1024)
	for {
		n, err := request.Body.Read(buf)
		if err != nil {
			fmt.Println("buf read ", err)
			break
		}
		_, err = create.Write(buf[:n])
	}
	_, err = rConn.Do("HSET", "MP_"+uploadID, "chkidx"+chunkIndex, 1)
	_, err = write.Write(util.NewRespMsg(0, "ok", nil).JSONBytes())
}
func CompleteUploadHandler(write http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		panic(err)
	}
	username := request.Form.Get("username")
	uploadID := request.Form.Get("uploadid")
	filehash := request.Form.Get("filehash")
	filesize := request.Form.Get("filesize")
	filename := request.Form.Get("filename")

	rConn := redis.RedisPool().Get()
	defer func() {
		err2 := rConn.Close()
		if err2 != nil {
			panic(err2)
		}
	}()
	values, err := redis2.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		_, err = write.Write(util.NewRespMsg(-1, "complete upload err"+err.Error(), nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(values); i += 2 {
		key := string(values[i].([]byte))
		value := string(values[i+1].([]byte))
		if key == "chunkcount" {
			totalCount, _ = strconv.Atoi(value)
		} else if strings.HasPrefix(key, "chkidx_") && value == "1" {
			chunkCount += 1
		}
	}
	if totalCount != chunkCount {
		_, err = write.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}

	atoi, _ := strconv.Atoi(filesize)
	db.OnfileUpdateFinish(filehash, filename, atoi, "")
	db.OnUserFileUploadFinished(username, filehash, filename, int64(atoi))
	_, err = write.Write(util.NewRespMsg(0, "ok", nil).JSONBytes())
}
