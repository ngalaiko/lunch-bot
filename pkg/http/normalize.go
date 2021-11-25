package http

import (
	"net/http"
)

func normalizePath(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		if lastChar := r.URL.Path[len(r.URL.Path)-1]; lastChar != '/' {
			r.URL.Path += "/"
		}
		next(w, r)
	}
}
