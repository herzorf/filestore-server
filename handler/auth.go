package handler

import (
	"net/http"
)

// HTTPIntercepter 请求拦截器
func HTTPIntercepter(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(write http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			write.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}
		username := request.Form.Get("username")
		token := request.Form.Get("token")
		if len(username) < 3 || !IsTokenValid(username, token) {
			write.WriteHeader(http.StatusForbidden)
			return
		}
		h(write, request)
	}
}
