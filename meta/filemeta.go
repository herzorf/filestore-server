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

var FileMetas map[string]FileMeta

func init() {
	FileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta 新增和更新fmeta
func UpdateFileMeta(fmeta FileMeta) {
	FileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB 更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnfileUpdateFinish(fmeta.FileSha1, fmeta.FileName, int(fmeta.FileSize), fmeta.Location)
}

// GetFileMeta 获取FileMetas里的元信息对象
func GetFileMeta(filesha1 string) FileMeta {
	return FileMetas[filesha1]
}

// RemoveFileMeta 删除文件元信息
func RemoveFileMeta(filesha1 string) {
	delete(FileMetas, filesha1)
}
