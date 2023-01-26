package main

import (
	"fmt"
	"github.com/herzorf/filestroe-server/route"
)

func main() {
	//cos.ConnectCos()
	router := route.Router()
	err := router.Run(":8080")
	if err != nil {
		fmt.Println("gin run err", err)
	}
	//http.HandleFunc("/api/file/meta", handler.GetFileMetaHandler)
	//http.HandleFunc("/api/file/download", handler.DownloadHandler)
	//http.HandleFunc("/api/file/update", handler.FileMetaUpdateHandler)
	//http.HandleFunc("/api/file/delete", handler.FileDeleteHandler)

	//http.HandleFunc("/api/file/fastUpload", handler.HTTPIntercepter(handler.TryFastUploadHandler))
	//http.HandleFunc("/api/file/mpupload/init", handler.HTTPIntercepter(handler.InitialMultipartUploadHandler))
	//http.HandleFunc("/api/file/mpupload/uppart", handler.HTTPIntercepter(handler.UploadPartHandler))
	//http.HandleFunc("/api/file/mpupload/complete", handler.HTTPIntercepter(handler.CompleteUploadHandler))

}
