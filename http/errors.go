package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	ErrStatusInternal       = "INTERNAL_SERVER_ERROR"
	ErrStatusUnauthorized   = "UNAUTHORIZED"
	ErrStatusNotFound       = "NOT_FOUND"
	ErrStatusBadRequest     = "MALFORMED_REQUEST"
	ErrStatusForbidden      = "PERMISSION_DENIED"
	ErrStatusConflict       = "CONFLICT"
	ErrStatusNotImplemented = "NOT_IMPLEMENTED"
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
		e.Logger.Error(e.String())
		e.Error = nil // e.Error may come from db, etc. So we hide this from the user
		w.WriteHeader(statusCode[e.Code])
		json.NewEncoder(w).Encode(e)
	}
}
