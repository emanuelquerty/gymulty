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
	"github.com/emanuelquerty/gymulty/http/middleware"
)

type ClassHandler struct {
	http.Handler
	store  domain.ClassStore
	logger *slog.Logger
}

func NewClassHandler(logger *slog.Logger, store domain.ClassStore) *ClassHandler {
	router := http.NewServeMux()

	handler := &ClassHandler{
		Handler: middleware.StripSlashes(router),
		store:   store,
		logger:  logger,
	}
	handler.registerRoutes(router)
	return handler
}

func (c *ClassHandler) registerRoutes(router *http.ServeMux) {
	router.Handle("POST /api/tenants/{tenantID}/classes", errorHandler(c.CreateClass))
	router.Handle("GET /api/tenants/{tenantID}/classes/{classID}", errorHandler(c.GetClassByID))
	router.Handle("DELETE /api/tenants/{tenantID}/classes/{classID}", errorHandler(c.DeleteClassByID))
}

func (c *ClassHandler) CreateClass(w http.ResponseWriter, r *http.Request) *appError {
	e := appError{Logger: c.logger}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, "Invalid tenant id", ErrStatusBadRequest)
	}

	var class domain.Class
	json.NewDecoder(r.Body).Decode(&class)

	class, err = c.store.CreateClass(r.Context(), tenantID, class)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.withContext(err, "Unknown tenant id", ErrStatusNotFound)
		}
		return e.withContext(err, "An internal server error ocurred. Please try again later", ErrStatusInternal)
	}

	resourceURI := fmt.Sprintf("%s://%s%s/%d", r.URL.Scheme, r.Host, r.URL.String(), class.ID)
	w.Header().Set("Location", resourceURI)

	w.WriteHeader(http.StatusCreated)
	res := Response[[]domain.Class]{Count: 1, Data: []domain.Class{class}}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return e.withContext(err, "An internal server error ocurred. Please try again later", ErrStatusInternal)
	}
	return nil
}

func (c *ClassHandler) GetClassByID(w http.ResponseWriter, r *http.Request) *appError {
	e := appError{Logger: c.logger}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, "Invalid tenant id", ErrStatusBadRequest)
	}

	classID, err := strconv.Atoi(r.PathValue("classID"))
	if err != nil {
		return e.withContext(err, "Invalid class id", ErrStatusBadRequest)
	}

	class, err := c.store.GetClassByID(r.Context(), tenantID, classID)
	fmt.Println("GET USER BY ID", class)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.withContext(err, "Unknown tenant or class id", ErrStatusNotFound)
		}
		return e.withContext(err, "An internal server error ocurred. Please try again later", ErrStatusInternal)
	}

	w.WriteHeader(http.StatusOK)
	res := Response[[]domain.Class]{Count: 1, Data: []domain.Class{class}}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return e.withContext(err, "An internal server error ocurred. Please try again later", ErrStatusInternal)
	}
	return nil
}

func (c *ClassHandler) DeleteClassByID(w http.ResponseWriter, r *http.Request) *appError {
	e := appError{Logger: c.logger}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, "Invalid tenant id", ErrStatusBadRequest)
	}

	classID, err := strconv.Atoi(r.PathValue("classID"))
	if err != nil {
		return e.withContext(err, "Invalid class id", ErrStatusBadRequest)
	}

	err = c.store.DeleteClassByID(r.Context(), tenantID, classID)
	if err != nil {
		return e.withContext(err, "An internal server error ocurred. Please try again later", ErrStatusInternal)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
