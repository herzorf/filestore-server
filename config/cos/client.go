package cos

import (
	"context"
	"fmt"
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
	// 1.永久密钥
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  SECRETID,
			SecretKey: SECRETKEY,
		},
	})

	s, _, err := client.Service.Get(context.Background())
	if err != nil {
		panic(err)
	}
	for _, b := range s.Buckets {
		fmt.Printf("%#v\n", b)
	}
	return client
}
