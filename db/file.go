package db

import (
	"fmt"
	"github.com/herzorf/filestroe-server/db/mysql"
)

// OnfileUpdateFinish 文件上传完成
func OnfileUpdateFinish(filehash string, fileName string, fileSize int, fileAddr string) bool {
	prepare, err := mysql.ConnectDB().Prepare("INSERT INTO file (file_sha1,file_name,file_size,file_addr,status) VALUES (?,?,?,?,1)")
	if err != nil {
		fmt.Println("prepare err", err)
		return false
	}
	defer func() {
		err = prepare.Close()
		if err != nil {
			fmt.Print("prepare close err")
			panic(err)
		}
	}()
	exec, err := prepare.Exec(filehash, fileName, fileSize, fileAddr)
	if err != nil {
		fmt.Println("insert err ", err)
		return false
	}

	if affected, err := exec.RowsAffected(); err == nil {
		if affected <= 0 {
			fmt.Printf("file with hash: %s hash been upload before\n", fileName)
		}
		return true
	}

	return false
}
