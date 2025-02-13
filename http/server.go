package http

import (
	"log/slog"
	"net/http"

	"github.com/emanuelquerty/gymulty/domain"
	"github.com/emanuelquerty/gymulty/http/middleware"
	"github.com/emanuelquerty/gymulty/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	router      http.Handler
	logger      *slog.Logger
	middlewares []middleware.Middleware
	store       domain.Store
}

func NewServer(pool *pgxpool.Pool, logger *slog.Logger) *Server {
	store := postgres.NewStore(pool)
	router := http.NewServeMux()

	server := &Server{
		router: router,
		logger: logger,
		store:  store,
	}

	server.registerRoutes(router)
	return server
}

func (s *Server) registerRoutes(router *http.ServeMux) {
	tenantHandler := NewTenantHandler(s.logger, s.store)
	userHandler := NewUserHandler(s.logger, s.store)
	classHandler := NewClassHandler(s.logger, s.store)

	router.Handle("/api/tenants/", tenantHandler)
	router.Handle("/api/tenants/{tenantID}/users/", userHandler)
	router.Handle("/api/tenants/{tenantID}/classes/", classHandler)
}

func (s *Server) Use(m middleware.Middleware) {
	s.middlewares = append(s.middlewares, m)
}

func (s *Server) registerGlobalMiddlewares() http.Handler {
	handler := s.router
	for _, m := range s.middlewares {
		handler = m(s.logger, handler)
	}
	return handler
}

func (s *Server) ListenAndServe(port string) error {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: s.registerGlobalMiddlewares(),
	}
	s.logger.Info("server is running", slog.String("port", port))
	return server.ListenAndServe()
}
