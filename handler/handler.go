package handler

import (
	"fmt"
	"net/http"
)

func UploadHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.Method)
}
