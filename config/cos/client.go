package cos

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

type BaseURL struct {
	BucketURL *url.URL
}

func ConnectCos() *cos.Client {
	u, _ := url.Parse("https://filestore-store-1304254779.cos.ap-shanghai.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  SECRETID,
			SecretKey: SECRETKEY,
		},
	})
	return client
}
