package cos

import (
	"context"
	"fmt"
	"io"
	"net/url"
)

func GetUploadObjectUrl(key string) *url.URL {
	objectUrl := ConnectCos().Object.GetObjectURL(key)
	fmt.Println(objectUrl)
	return objectUrl
}

func PutFileObject(file io.Reader, name string) error {
	_, err := ConnectCos().Object.Put(context.Background(), name, file, nil)
	if err != nil {
		return err
	}
	return nil
}
