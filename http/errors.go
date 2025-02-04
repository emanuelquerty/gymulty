package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	ErrInternal       = "INTERNAL_SERVER_ERROR"
	ErrUnauthorized   = "UNAUTHORIZED"
	ErrNotFound       = "NOT_FOUND"
	ErrBadRequest     = "MALFORMED_REQUEST"
	ErrForbidden      = "PERMISSION_DENIED"
	ErrConflict       = "CONFLICT"
	ErrNotImplemented = "NOT_IMPLEMENTED"
)

var statusCode = map[string]int{
	ErrInternal:       http.StatusInternalServerError,
	ErrUnauthorized:   http.StatusUnauthorized,
	ErrNotFound:       http.StatusNotFound,
	ErrBadRequest:     http.StatusBadRequest,
	ErrForbidden:      http.StatusForbidden,
	ErrConflict:       http.StatusConflict,
	ErrNotImplemented: http.StatusNotImplemented,
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

type errorHandler func(w http.ResponseWriter, r *http.Request) *appError

func (fn errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not error
		e.Logger.Error(e.Error.Error())
		e.Error = nil // e.Error may come from db, etc. So we hide this from the user
		w.WriteHeader(statusCode[e.Code])
		json.NewEncoder(w).Encode(e)
	}
}
