package db

import (
	"fmt"
	"github.com/herzorf/filestroe-server/db/mysql"
	"time"
)

type UserFile struct {
	Username    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

// OnUserFileUploadFinished 更新用户文件表信息
func OnUserFileUploadFinished(username, fileHash, fileName string, fileSize int64) bool {
	stmt, err := mysql.ConnectDB().Prepare("INSERT INTO user_file (user_name,file_sha1,file_size,file_name,upload_at) VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Println("prepare err", err)
		return false
	}
	defer func() {
		err2 := stmt.Close()
		if err2 != nil {
			panic(err2)
		}
	}()
	_, err = stmt.Exec(fileName, fileHash, fileSize, fileName, time.Now())
	if err != nil {
		fmt.Println("insert err", err)
		return false
	}
	return true
}
