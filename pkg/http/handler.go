package http

import (
	"net/http"
)

type handler struct {
	mux *http.ServeMux

	defaultMiddlewares []middleware
}

func newHandler() *handler {
	return &handler{
		mux:                http.NewServeMux(),
		defaultMiddlewares: []middleware{normalizePath, accessLogs},
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

func (h *handler) POST(path string, handler http.HandlerFunc, mm ...middleware) {
	for _, m := range append(mm, ensureMethod(http.MethodPost), ensurePath(path)) {
		handler = m(handler)
	}
	h.mux.HandleFunc(path, handler)
}
