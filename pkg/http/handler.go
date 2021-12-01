package http

import (
	"net/http"
)

type middleware func(http.HandlerFunc) http.HandlerFunc

type handler struct {
	mux *http.ServeMux

	defaultMiddlewares []middleware
}

func newHandler(mm ...middleware) *handler {
	return &handler{
		mux:                http.NewServeMux(),
		defaultMiddlewares: mm,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := h.mux.ServeHTTP
	for _, m := range h.defaultMiddlewares {
		handler = m(handler)
	}
	handler(w, r)
}

func (h *handler) GET(path string, handler http.HandlerFunc, mm ...middleware) {
	for _, m := range append(mm, ensureMethod(http.MethodGet), ensurePath(path)) {
		handler = m(handler)
	}
	h.mux.HandleFunc(path, handler)
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
