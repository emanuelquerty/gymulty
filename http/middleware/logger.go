package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func Logger(logger *slog.Logger, fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWritter{w: w}

		then := time.Now()
		fn.ServeHTTP(lrw, r)
		duration := time.Since(then).Milliseconds()

		requestAttrs := slog.Group("request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.String()),
			slog.String("proto", r.Proto),
		)

		responseAttrs := slog.Group("response",
			slog.Int("status_code", lrw.statusCode),
			slog.Int("size", lrw.bytesCount),
			slog.String("duration", fmt.Sprintf("%d ms", duration)),
		)

		logger.Info("REQUEST", requestAttrs)
		logger.Info("RESPONSE", responseAttrs)
	})
}
