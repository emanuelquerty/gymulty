package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/emanuelquerty/gymulty/http/middleware"
)

const (
	ErrMsgInternal          = "An internal server error ocurred. Please try again later"
	ErrMsgInvalidResourceID = "Invalid resource id"
)

const (
	ErrStatusInternal       = "internal_server_error"
	ErrStatusUnauthorized   = "unauthorized"
	ErrStatusNotFound       = "not_found"
	ErrStatusBadRequest     = "malformed_request"
	ErrStatusForbidden      = "permission_denied"
	ErrStatusConflict       = "conflict"
	ErrStatusNotImplemented = "not_implemented"
)

var statusCode = map[string]int{
	ErrStatusInternal:       http.StatusInternalServerError,
	ErrStatusUnauthorized:   http.StatusUnauthorized,
	ErrStatusNotFound:       http.StatusNotFound,
	ErrStatusBadRequest:     http.StatusBadRequest,
	ErrStatusForbidden:      http.StatusForbidden,
	ErrStatusConflict:       http.StatusConflict,
	ErrStatusNotImplemented: http.StatusNotImplemented,
}

type appError struct {
	Error   error        `json:"error,omitempty"  bson:"error"`
	Code    string       `json:"code,omitempty"  bson:"code"`
	Message string       `json:"message,omitempty"  bson:"message"`
	Logger  *slog.Logger `json:"-"  bson:"-"`
}

func (e *appError) withContext(err error, msg string, statusText string) *appError {
	e.Error = err
	e.Message = msg
	e.Code = statusText
	return e
}

func (e *appError) String() string {
	return e.Error.Error()
}

type errorHandler func(w http.ResponseWriter, r *http.Request) *appError

func (fn errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not error
		reqID := middleware.GetRequestID(r.Context(), e.Logger)
		e.Logger.Error(e.Message, slog.String("request_id", reqID), slog.String("error", e.String()))
		e.Error = nil // e.Error may come from db, etc. So we hide this from the user
		w.WriteHeader(statusCode[e.Code])
		json.NewEncoder(w).Encode(e)
	}
}
