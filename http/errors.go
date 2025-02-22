package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/emanuelquerty/gymulty/http/middleware"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrMsgInternal          = "An internal server error ocurred. Please try again later"
	ErrMsgInvalidResourceID = "Invalid resource id"
	ErrMsgNotFound          = "The resource with specified id was not found"
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

var constraintErrors = map[string]string{
	"tenants_subdomain_key": "Subdomain already exists",
	"tenants_status_check":  "Invalid value for status",
	"users_email_key":       "Email already exists",
	"users_role_check":      "Invalid value for role",
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

	if errors.Is(err, sql.ErrNoRows) {
		e.Message = ErrMsgNotFound
		e.Code = ErrStatusNotFound
	}

	if dbError, ok := err.(*pgconn.PgError); ok {
		if msg, exists := constraintErrors[dbError.ConstraintName]; exists {
			e.Message = msg
		}
		switch dbError.Code {
		case "23505":
			e.Code = ErrStatusConflict
		case "23514":
			e.Code = ErrStatusBadRequest
		}
	}
	return e
}

func (e *appError) String() string {
	return e.Error.Error()
}

type errorHandler func(w http.ResponseWriter, r *http.Request) *appError

func (fn errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		reqID := middleware.GetRequestID(r.Context(), e.Logger)
		e.Logger.Error(e.Message,
			slog.String("request_id", reqID),
			slog.String("error", e.String()),
		)
		e.Error = nil // e.Error is for logging only.
		w.WriteHeader(statusCode[e.Code])
		json.NewEncoder(w).Encode(e)
	}
}
