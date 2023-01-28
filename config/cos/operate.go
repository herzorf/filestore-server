package cos

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/url"
)

func GetUploadObjectUrl(key string) *url.URL {
	objectUrl := ConnectCos().Object.GetObjectURL(key)
	return objectUrl
}

func PutFileObject(file io.Reader, name string, contentType string) error {
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:        contentType,
			ContentDisposition: "attachment",
		},
	}
	_, err := ConnectCos().Object.Put(context.Background(), name, file, opt)
	if err != nil {
		return err
	}
	return nil
}
