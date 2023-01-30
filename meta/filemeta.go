package meta

import "github.com/herzorf/filestroe-server/db"

// FileMeta 文件元信息
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnfileUpdateFinish(fmeta.FileSha1, fmeta.FileName, int(fmeta.FileSize), fmeta.Location)
}
