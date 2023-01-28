package db

import (
	"fmt"
	"github.com/herzorf/filestroe-server/db/mysql"
	"log"
	"time"
)

type UserFile struct {
	Username    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
	Location    string
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
	_, err = stmt.Exec(username, fileHash, fileSize, fileName, time.Now())
	if err != nil {
		fmt.Println("insert err", err)
		return false
	}
	return true
}

func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	var userFiles []UserFile
	stmt, err := mysql.ConnectDB().Prepare("SELECT file_sha1,file_name,file_size,upload_at,last_update FROM user_file WHERE user_name = ? LIMIT ?")
	if err != nil {
		fmt.Println("prepare err", err)
		return userFiles, err
	}
	defer func() {
		err2 := stmt.Close()
		if err2 != nil {
			panic(err2)
		}
	}()
	rows, err := stmt.Query(username, limit)
	if err != nil {
		fmt.Println("prepare err", err)
		return userFiles, err
	}
	for rows.Next() {
		userFile := UserFile{}
		err := rows.Scan(&userFile.FileHash, &userFile.FileName, &userFile.FileSize, &userFile.UploadAt, &userFile.LastUpdated)
		if err != nil {
			log.Fatal(err)
		}
		userFile.Username = username
		userFiles = append(userFiles, userFile)
	}
	return userFiles, nil
}

func DeleteUserFileMetas(hash string) error {
	prepare, err := mysql.ConnectDB().Prepare("DELETE FROM user_file WHERE file_sha1 = ?")
	if err != nil {
		return err
	}
	defer func() {
		err2 := prepare.Close()
		if err2 != nil {
			log.Println("prepare close err", err)
		}
	}()
	exec, err := prepare.Exec(hash)
	if err != nil {
		return err
	}
	log.Println(exec.RowsAffected())
	return nil
}
