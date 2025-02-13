package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/emanuelquerty/gymulty/domain"
)

type ClassHandler struct {
	http.Handler
	store  domain.ClassStore
	logger *slog.Logger
}

func NewClassHandler(logger *slog.Logger, store domain.ClassStore) *ClassHandler {
	router := http.NewServeMux()

	handler := &ClassHandler{
		Handler: router,
		store:   store,
		logger:  logger,
	}
	handler.registerRoutes(router)
	return handler
}

func (c *ClassHandler) registerRoutes(router *http.ServeMux) {
	router.Handle("GET /api/tenants/{tenantID}/classes", errorHandler(c.CreateClass))
}

func (c *ClassHandler) CreateClass(w http.ResponseWriter, r *http.Request) *appError {
	e := appError{Logger: c.logger}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, "Invalid tenant id", ErrBadRequest)
	}

	var class domain.Class
	json.NewDecoder(r.Body).Decode(&class)

	class, err = c.store.CreateClass(r.Context(), tenantID, class)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.withContext(err, "A tenant with specified id was not found", ErrNotFound)
		}
		return e.withContext(err, "An internal server error ocurred. Please try again later", ErrInternal)
	}

	resourceURI := fmt.Sprintf("%s://%s%s/%d", r.URL.Scheme, r.Host, r.URL.String(), class.ID)
	w.Header().Set("Location", resourceURI)

	err = json.NewEncoder(w).Encode(class)
	if err != nil {
		return e.withContext(err, "An internal server error ocurred. Please try again later", ErrInternal)
	}
	return nil
}
