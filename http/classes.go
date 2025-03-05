package http

import (
	"encoding/json"
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
	router.Handle("GET /api/tenants/{tenantID}/classes", errorHandler(c.GetAllClasses))
}

// CreateClass godoc
//
//	@Summary		Create new class
//	@Description	Creates a new class
//	@Tags			classes
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		int			       true	      "The tenant ID"
//	@Param			class		body		domain.Class	   true       "The request body"
//	@Success		200         {object}	Response[[]domain.Class]
//	@Failure		500			{object}	appError
//	@Failure		409			{object}	appError
//	@Failure		404			{object}	appError
//	@Failure		400			{object}	appError
//	@Router			/api/tenants/{tenantID}/classes   [post]
func (c *ClassHandler) CreateClass(w http.ResponseWriter, r *http.Request) *appError {
	e := appError{Logger: c.logger}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	var class domain.Class
	json.NewDecoder(r.Body).Decode(&class)

	class, err = c.store.CreateClass(r.Context(), tenantID, class)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}

	resourceURI := fmt.Sprintf("%s://%s%s/%d", r.URL.Scheme, r.Host, r.URL.String(), class.ID)
	w.Header().Set("Location", resourceURI)

	w.WriteHeader(http.StatusCreated)
	res := Response[[]domain.Class]{Count: 1, Data: []domain.Class{class}}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}
	return nil
}

// GetClassByID godoc
//
//	@Summary		Get class by id
//	@Description	Fetches a class by its id
//	@Tags			classes
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		int			       true	      "The tenant ID"
//	@Param			classID	    path		int			       true	      "The class ID"
//	@Success		200         {object}	Response[[]domain.Class]
//	@Failure		500			{object}	appError
//	@Failure		404			{object}	appError
//	@Failure		400			{object}	appError
//	@Router			/api/tenants/{tenantID}/classes/{classID}   [get]
func (c *ClassHandler) GetClassByID(w http.ResponseWriter, r *http.Request) *appError {
	e := appError{Logger: c.logger}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	classID, err := strconv.Atoi(r.PathValue("classID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	class, err := c.store.GetClassByID(r.Context(), tenantID, classID)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}

	w.WriteHeader(http.StatusOK)
	res := Response[[]domain.Class]{Count: 1, Data: []domain.Class{class}}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}
	return nil
}

// DeleteClassByID godoc
//
//	@Summary		Delete class by id
//	@Description	Delete a class by id
//	@Tags			classes
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		int			       true	      "The tenant ID"
//	@Param			classID	    path		int			       true	      "The class ID"
//	@Success		204
//	@Failure		500			{object}	appError
//	@Failure		404			{object}	appError
//	@Failure		400			{object}	appError
//	@Router			/api/tenants/{tenantID}/classes/{classID}   [delete]
func (c *ClassHandler) DeleteClassByID(w http.ResponseWriter, r *http.Request) *appError {
	e := appError{Logger: c.logger}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	classID, err := strconv.Atoi(r.PathValue("classID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	err = c.store.DeleteClassByID(r.Context(), tenantID, classID)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// GetAll classes godoc
//
//	@Summary		Fetch all classes
//	@Description	Fetches all classes for a given tenant
//	@Tags			classes
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		int			       true	      "The tenant ID"
//	@Success		200         {object}	Response[[]domain.Class]
//	@Failure		500			{object}	appError
//	@Failure		400			{object}	appError
//	@Router			/api/tenants/{tenantID}/classes   [get]
func (c *ClassHandler) GetAllClasses(w http.ResponseWriter, r *http.Request) *appError {
	e := appError{Logger: c.logger}

	tenantID, err := strconv.Atoi(r.PathValue("tenantID"))
	if err != nil {
		return e.withContext(err, ErrMsgInvalidResourceID, ErrStatusBadRequest)
	}

	classes, err := c.store.GetAllClasses(r.Context(), tenantID)
	if err != nil {
		return e.withContext(err, ErrMsgInternal, ErrStatusInternal)
	}
	res := Response[[]domain.Class]{
		Count: len(classes),
		Data:  classes,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return nil
}
