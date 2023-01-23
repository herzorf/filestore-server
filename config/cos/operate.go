package cos

import (
	"context"
	"fmt"
	"net/url"
	"os"
)

func GetUploadObjectUrl(key string) *url.URL {
	objectUrl := ConnectCos().Object.GetObjectURL(key)
	fmt.Println(objectUrl)
	return objectUrl
}

func PutFileObject(file *os.File, name string) error {

	_, err := ConnectCos().Object.Put(context.Background(), name, file, nil)
	if err != nil {
		return err
	}
	return nil
}
