package middleware

import (
	"log/slog"
	"net/http"
)

type Middleware func(logger *slog.Logger, handler http.Handler) http.Handler

type loggingResponseWritter struct {
	w          http.ResponseWriter
	statusCode int
	bytesCount int
}

func (lrw *loggingResponseWritter) Header() http.Header {
	return lrw.w.Header()
}

func (lrw *loggingResponseWritter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.w.WriteHeader(statusCode)
}

func (lrw *loggingResponseWritter) Write(b []byte) (int, error) {
	lrw.bytesCount += len(b)
	return lrw.w.Write(b)
}
