package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func Logger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := GetRequestID(r.Context(), logger)
		lrw := &loggingResponseWritter{w: w}

		then := time.Now()
		next.ServeHTTP(lrw, r)
		duration := time.Since(then).Milliseconds()

		if lrw.statusCode == http.StatusMovedPermanently {
			return
		}

		requestAttrs := slog.Group("request",
			slog.String("id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.String()),
			slog.String("proto", r.Proto),
		)
		responseAttrs := slog.Group("response",
			slog.Int("status_code", lrw.statusCode),
			slog.Int("size", lrw.bytesCount),
			slog.String("duration", fmt.Sprintf("%d ms", duration)),
		)

		logger.Info("request", requestAttrs)
		logger.Info("response", responseAttrs)
	})
}
