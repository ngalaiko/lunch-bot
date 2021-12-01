package http

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"time"
)

type middleware func(http.HandlerFunc) http.HandlerFunc

type loggingResponseWriter struct {
	http.ResponseWriter

	StatusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return lrw.ResponseWriter.(http.Hijacker).Hijack()
}

func accessLogs(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wl := newLoggingResponseWriter(w)

		next(wl, r.Clone(r.Context()))

		log.Printf("[INFO] %s %s %d %s", r.Method, r.URL, wl.StatusCode, time.Since(start))
	}
}

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

func ensureMethod(method string) middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			next(w, r)
		}
	}
}

func ensurePath(path string) middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != path {
				http.NotFound(w, r)
				return
			}
			next(w, r)
		}
	}
}
