package db

import (
	"database/sql"
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

type TableFile struct {
	Filehash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// OnGetFileMeta 从数据库中获取元信息
func OnGetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mysql.ConnectDB().Prepare("SELECT file_sha1,file_addr,file_name,file_size FROM file WHERE file_sha1 = ? AND status = 1 LIMIT 1")
	if err != nil {
		fmt.Println("mysql prepare err", err)
		return nil, err
	}
	defer func() {
		err2 := stmt.Close()
		if err2 != nil {
			fmt.Println("stmt close err", err)
		}
	}()
	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.Filehash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		fmt.Println("数据库查询错误", err)
		return nil, err
	}
	return &tfile, nil
}
