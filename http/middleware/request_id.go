package middleware

import (
	"context"
	"encoding/base64"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type requestIdCtxKeyType string

const reqIdCtxKey requestIdCtxKeyType = "requestID"

func AddRequestID(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()

		b64ID := base64.StdEncoding.EncodeToString(id[:])
		b64ID = strings.NewReplacer("+", "", "==", "", "/", "").Replace(b64ID)

		ctx := withRequestID(r.Context(), b64ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withRequestID(ctx context.Context, requestID string) context.Context {
	ctx = context.WithValue(ctx, reqIdCtxKey, requestID)
	return ctx
}

func GetRequestID(ctx context.Context, logger *slog.Logger) string {
	reqID, ok := ctx.Value(reqIdCtxKey).(string)
	if !ok {
		logger.Error("Retrieving request id", slog.String("error", "request id was not found"))
		return ""
	}
	return string(reqID)
}
