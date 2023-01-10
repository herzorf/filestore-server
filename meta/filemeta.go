package meta

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

// GetFileMeta 获取FileMetas里的元信息对象
func GetFileMeta(filesha1 string) FileMeta {
	return FileMetas[filesha1]
}
